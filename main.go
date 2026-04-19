package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	keyFocus = iota
	valueFocus
	tableFocus
	focusCount
)

type keyMap struct {
	Next       key.Binding
	Previous   key.Binding
	Save       key.Binding
	FocusTable key.Binding
	Delete     key.Binding
	Refresh    key.Binding
	Quit       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Next, k.Save, k.Delete, k.Refresh, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Previous, k.Save, k.FocusTable},
		{k.Delete, k.Refresh, k.Quit},
	}
}

var keys = keyMap{
	Next:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "focus")),
	Previous:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
	Save:       key.NewBinding(key.WithKeys("enter", "ctrl+s"), key.WithHelp("enter", "save/edit")),
	FocusTable: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "table")),
	Delete:     key.NewBinding(key.WithKeys("d", "delete"), key.WithHelp("d", "delete")),
	Refresh:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
}

type listMsg struct {
	rows      []table.Row
	rawValues map[string]string
	err       error
}

type setMsg struct {
	key string
	err error
}

type deleteMsg struct {
	key string
	err error
}

type model struct {
	keyInput   textinput.Model
	valueInput textinput.Model
	table      table.Model
	help       help.Model
	spinner    spinner.Model

	focus   int
	width   int
	height  int
	loading bool
	status  string
	raw     map[string]string

	keyColumnWidth   int
	valueColumnWidth int
}

var (
	purple      = lipgloss.Color("99")
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	tableFrame  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Padding(1, 1)
)

func main() {
	if _, err := exec.LookPath("skate"); err != nil {
		fmt.Fprintln(os.Stderr, "skate was not found in PATH")
		os.Exit(1)
	}

	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newModel() model {
	keyInput := textinput.New()
	keyInput.Placeholder = "key"
	keyInput.Prompt = ""

	valueInput := textinput.New()
	valueInput.Placeholder = "value"
	valueInput.Prompt = ""

	columns := []table.Column{
		{Title: "Key", Width: 24},
		{Title: "Value", Width: 48},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithHeight(10),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.Bold(true).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true)
	styles.Selected = styles.Selected.Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62")).Bold(false)
	t.SetStyles(styles)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = statusStyle

	m := model{
		keyInput:         keyInput,
		valueInput:       valueInput,
		table:            t,
		help:             help.New(),
		spinner:          s,
		loading:          true,
		status:           "loading keys",
		raw:              map[string]string{},
		keyColumnWidth:   columns[0].Width,
		valueColumnWidth: columns[1].Width,
	}
	m.setFocus(0)
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(loadRowsCmd(), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()
		return m, nil

	case spinner.TickMsg:
		if !m.loading {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit) && (msg.String() == "ctrl+c" || m.focus == tableFocus):
			return m, tea.Quit
		case key.Matches(msg, keys.Next):
			m.setFocus((m.focus + 1) % focusCount)
			return m, nil
		case key.Matches(msg, keys.Previous):
			m.setFocus((m.focus + focusCount - 1) % focusCount)
			return m, nil
		case key.Matches(msg, keys.Save):
			switch m.focus {
			case keyFocus, valueFocus:
				return m.save()
			case tableFocus:
				m.fillInputsFromSelection()
				m.setFocus(valueFocus)
				return m, nil
			}
		case key.Matches(msg, keys.FocusTable):
			m.setFocus(tableFocus)
			return m, nil
		}

		if m.focus == tableFocus {
			switch {
			case key.Matches(msg, keys.Refresh):
				m.loading = true
				m.status = "refreshing"
				return m, tea.Batch(loadRowsCmd(), m.spinner.Tick)
			case key.Matches(msg, keys.Delete):
				row := m.table.SelectedRow()
				if len(row) == 0 {
					m.status = "no key selected"
					return m, nil
				}
				key := row[0]
				m.loading = true
				m.status = "deleting " + key
				return m, tea.Batch(deleteCmd(key), m.spinner.Tick)
			}
		}

	case listMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "list failed: " + msg.err.Error()
			return m, nil
		}
		m.table.SetRows(msg.rows)
		m.raw = msg.rawValues
		m.resize()
		m.status = fmt.Sprintf("%d keys", len(msg.rows))
		return m, nil

	case setMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "save failed: " + msg.err.Error()
			return m, nil
		}
		m.status = "saved " + msg.key
		m.keyInput.SetValue("")
		m.valueInput.SetValue("")
		m.loading = true
		return m, tea.Batch(loadRowsCmd(), m.spinner.Tick)

	case deleteMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "delete failed: " + msg.err.Error()
			return m, nil
		}
		if m.keyInput.Value() == msg.key {
			m.keyInput.SetValue("")
			m.valueInput.SetValue("")
		}
		m.status = "deleted " + msg.key
		m.loading = true
		return m, tea.Batch(loadRowsCmd(), m.spinner.Tick)
	}

	switch m.focus {
	case keyFocus:
		var cmd tea.Cmd
		m.keyInput, cmd = m.keyInput.Update(msg)
		m.resize()
		cmds = append(cmds, cmd)
	case valueFocus:
		var cmd tea.Cmd
		m.valueInput, cmd = m.valueInput.Update(msg)
		m.resize()
		cmds = append(cmds, cmd)
	case tableFocus:
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	status := statusStyle.Render(m.status)
	if strings.Contains(m.status, "failed:") || strings.Contains(m.status, "required") {
		status = errorStyle.Render(m.status)
	}
	if m.loading {
		status = statusStyle.Render(m.spinner.View() + " " + m.status)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		lipgloss.NewStyle().MarginLeft(2).Render(m.inputRowView()),
		"",
		tableFrame.Render(m.table.View()),
		"",
		status,
		m.help.View(keys),
	)
}

func (m *model) setFocus(next int) {
	m.focus = next

	m.keyInput.Blur()
	m.valueInput.Blur()
	m.table.Blur()

	switch next {
	case keyFocus:
		m.keyInput.Focus()
	case valueFocus:
		m.valueInput.Focus()
	case tableFocus:
		m.table.Focus()
	}

	m.applyInputStyles()
}

func (m model) save() (tea.Model, tea.Cmd) {
	key := strings.TrimSpace(m.keyInput.Value())
	if key == "" {
		m.status = "key is required"
		return m, nil
	}
	m.loading = true
	m.status = "saving " + key
	return m, tea.Batch(setCmd(key, m.valueInput.Value()), m.spinner.Tick)
}

func (m *model) resize() {
	width := m.width
	if width <= 0 {
		width = 80
	}

	maxWidth := max(52, width-tableFrame.GetHorizontalFrameSize())
	contentWidth := min(maxWidth, m.desiredWidth())
	m.help.Width = contentWidth
	m.table.SetWidth(contentWidth)

	keyWidth := min(max(12, m.desiredKeyWidth()), max(18, contentWidth/3))
	valueWidth := max(12, contentWidth-keyWidth-5)
	m.keyColumnWidth = keyWidth
	m.valueColumnWidth = valueWidth
	m.keyInput.Width = keyWidth
	m.valueInput.Width = valueWidth
	m.table.SetColumns([]table.Column{
		{Title: "Key", Width: keyWidth},
		{Title: "Value", Width: valueWidth},
	})

	tableHeight := min(max(5, m.height-10), 9)
	if tableHeight < 5 {
		tableHeight = 5
	}
	m.table.SetHeight(tableHeight)
}

func (m model) desiredWidth() int {
	keyWidth := min(max(12, m.desiredKeyWidth()), 24)
	valueWidth := max(12, m.desiredValueWidth())

	return max(52, keyWidth+valueWidth+5)
}

func (m model) desiredKeyWidth() int {
	width := max(lipgloss.Width("Key"), lipgloss.Width(m.keyInput.Value()))
	for _, row := range m.table.Rows() {
		if len(row) > 0 {
			width = max(width, lipgloss.Width(row[0]))
		}
	}
	return width
}

func (m model) desiredValueWidth() int {
	width := max(lipgloss.Width("Value"), lipgloss.Width(m.valueInput.Value()))
	for _, row := range m.table.Rows() {
		if len(row) > 1 {
			width = max(width, lipgloss.Width(row[1]))
		}
	}
	return width
}

func (m model) inputRowView() string {
	cellStyle := lipgloss.NewStyle().Padding(0, 1)
	innerStyle := func(width int) lipgloss.Style {
		return lipgloss.NewStyle().Width(width).MaxWidth(width).Inline(true)
	}

	keyCell := cellStyle.Render(innerStyle(m.keyColumnWidth).Render(m.keyInput.View()))
	valueCell := cellStyle.Render(innerStyle(m.valueColumnWidth).Render(m.valueInput.View()))

	return lipgloss.JoinHorizontal(lipgloss.Top, keyCell, valueCell)
}

func (m *model) applyInputStyles() {
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	active := lipgloss.NewStyle().Foreground(purple)
	plain := lipgloss.NewStyle()

	m.keyInput.TextStyle = plain
	m.keyInput.PlaceholderStyle = muted
	m.keyInput.Cursor.Style = plain.Reverse(true)
	m.keyInput.Cursor.TextStyle = plain
	m.valueInput.TextStyle = plain
	m.valueInput.PlaceholderStyle = muted
	m.valueInput.Cursor.Style = plain.Reverse(true)
	m.valueInput.Cursor.TextStyle = plain

	if m.focus == keyFocus {
		m.keyInput.TextStyle = active
		m.keyInput.PlaceholderStyle = active
		m.keyInput.Cursor.Style = active.Reverse(true)
		m.keyInput.Cursor.TextStyle = active
	}

	if m.focus == valueFocus {
		m.valueInput.TextStyle = active
		m.valueInput.PlaceholderStyle = active
		m.valueInput.Cursor.Style = active.Reverse(true)
		m.valueInput.Cursor.TextStyle = active
	}
}

func (m *model) fillInputsFromSelection() {
	row := m.table.SelectedRow()
	if len(row) < 1 {
		return
	}
	m.keyInput.SetValue(row[0])
	m.valueInput.SetValue(m.raw[row[0]])
	m.resize()
}

func loadRowsCmd() tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("skate", "list", "--keys-only").CombinedOutput()
		if err != nil {
			return listMsg{err: commandErr(err, out)}
		}

		keys := parseKeys(out)
		rows := make([]table.Row, 0, len(keys))
		rawValues := make(map[string]string, len(keys))
		for _, key := range keys {
			valueOut, err := exec.Command("skate", "get", key).CombinedOutput()
			if err != nil {
				return listMsg{err: commandErr(err, valueOut)}
			}
			value := string(valueOut)
			rawValues[key] = value
			rows = append(rows, table.Row{key, displayValue(value)})
		}

		return listMsg{rows: rows, rawValues: rawValues}
	}
}

func setCmd(key string, value string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("skate", "set", key)
		cmd.Stdin = strings.NewReader(value)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return setMsg{key: key, err: commandErr(err, out)}
		}
		return setMsg{key: key}
	}
}

func deleteCmd(key string) tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("skate", "delete", key).CombinedOutput()
		if err != nil {
			return deleteMsg{key: key, err: commandErr(err, out)}
		}
		return deleteMsg{key: key}
	}
}

func parseKeys(out []byte) []string {
	keys := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if line == "" {
			continue
		}
		keys = append(keys, line)
	}
	return keys
}

func displayValue(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	return strings.ReplaceAll(value, "\n", "\\n")
}

func commandErr(err error, out []byte) error {
	msg := strings.TrimSpace(string(out))
	if msg == "" {
		return err
	}
	return fmt.Errorf("%w: %s", err, msg)
}
