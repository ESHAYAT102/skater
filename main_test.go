package main

import (
	"testing"

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
