package config

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestParseConsulValue_Base64(t *testing.T) {
	m := map[string]any{"a": 1}
	b, _ := json.Marshal(m)
	raw := base64.StdEncoding.EncodeToString(b)
	got, ok := parseConsulValue(raw)
	if !ok || got["a"].(float64) != 1 {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestParseConsulValue_PlainJSON(t *testing.T) {
	raw := `{"b":2}`
	got, ok := parseConsulValue(raw)
	if !ok || got["b"].(float64) != 2 {
		t.Fatalf("unexpected: %#v", got)
	}
}
