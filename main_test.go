package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func TestParseKeys(t *testing.T) {
	keys := parseKeys([]byte("alpha\nbeta\n"))

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "alpha" {
		t.Fatalf("unexpected first key: %q", keys[0])
	}
	if keys[1] != "beta" {
		t.Fatalf("unexpected second key: %q", keys[1])
	}
}

func TestParseKeysSkipsEmptyLines(t *testing.T) {
	keys := parseKeys([]byte("alpha\n\nbeta\n"))

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "alpha" || keys[1] != "beta" {
		t.Fatalf("unexpected keys: %#v", keys)
	}
}

func TestDisplayValueEscapesNewlines(t *testing.T) {
	got := displayValue("one\ntwo\r\nthree")

	if got != `one\ntwo\nthree` {
		t.Fatalf("unexpected display value: %q", got)
	}
}

func TestResizeColumnsFillTableWidth(t *testing.T) {
	m := newModel()
	m.width = 80
	m.height = 24
	m.resize()

	got := m.keyColumnWidth + m.valueColumnWidth + tableCellHorizontalFrameSize
	if got != m.table.Width() {
		t.Fatalf("column widths plus cell padding = %d, table width = %d", got, m.table.Width())
	}
	if got := lipgloss.Width(m.inputRowView()); got != m.table.Width() {
		t.Fatalf("input row width = %d, table width = %d", got, m.table.Width())
	}
}

func TestResizeColumnsFitContentBeforeTerminalMax(t *testing.T) {
	m := newModel()
	key := strings.Repeat("k", 40)
	value := "short"
	m.table.SetRows([]table.Row{{key, value}})
	m.width = 120
	m.height = 24
	m.resize()

	maxWidth := m.width - tableFrame.GetHorizontalFrameSize()
	if m.table.Width() >= maxWidth {
		t.Fatalf("table width = %d, max width = %d; expected content-sized table", m.table.Width(), maxWidth)
	}
	if m.keyColumnWidth < lipgloss.Width(key) {
		t.Fatalf("key column width = %d, key width = %d", m.keyColumnWidth, lipgloss.Width(key))
	}
	if m.valueColumnWidth < lipgloss.Width(value) {
		t.Fatalf("value column width = %d, value width = %d", m.valueColumnWidth, lipgloss.Width(value))
	}
}

func TestResizeColumnsUseTerminalMaxWhenContentTooWide(t *testing.T) {
	m := newModel()
	m.table.SetRows([]table.Row{{"key", strings.Repeat("v", 120)}})
	m.width = 80
	m.height = 24
	m.resize()

	maxWidth := m.width - tableFrame.GetHorizontalFrameSize()
	if m.table.Width() != maxWidth {
		t.Fatalf("table width = %d, max width = %d", m.table.Width(), maxWidth)
	}
	if got := m.keyColumnWidth + m.valueColumnWidth + tableCellHorizontalFrameSize; got != m.table.Width() {
		t.Fatalf("column widths plus cell padding = %d, table width = %d", got, m.table.Width())
	}
}
