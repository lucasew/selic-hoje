package handler

import (
	"strings"
	"testing"
)

func TestParseMetaSelic(t *testing.T) {
	const in = `[{"data":"05/08/2026","valor":"14.25"}]`
	got, err := parseMetaSelic([]byte(in))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "14.25" {
		t.Errorf("expected 14.25 got %q", got)
	}
}

func TestParseMetaSelicEmpty(t *testing.T) {
	if _, err := parseMetaSelic([]byte(`[]`)); err == nil {
		t.Fatal("expected error for empty series")
	}
	if _, err := parseMetaSelic([]byte(`[{"data":"01/01/2026","valor":""}]`)); err == nil {
		t.Fatal("expected error for empty valor")
	}
}

func TestParseMetaSelicInvalidJSON(t *testing.T) {
	_, err := parseMetaSelic([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "decode") {
		t.Errorf("error should mention decode: %v", err)
	}
}
