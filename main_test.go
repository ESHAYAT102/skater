package main

import "testing"

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
