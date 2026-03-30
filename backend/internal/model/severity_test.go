package model_test

import (
	"testing"

	"github.com/akaitigo/review-gym/internal/model"
)

func TestSeverityValid(t *testing.T) {
	tests := []struct {
		severity model.Severity
		want     bool
	}{
		{model.SeverityCritical, true},
		{model.SeverityMajor, true},
		{model.SeverityMinor, true},
		{model.SeverityInfo, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.severity), func(t *testing.T) {
			if got := tt.severity.Valid(); got != tt.want {
				t.Errorf("Severity(%q).Valid() = %v, want %v", tt.severity, got, tt.want)
			}
		})
	}
}

func TestAllSeverities(t *testing.T) {
	sevs := model.AllSeverities()
	if len(sevs) != 4 {
		t.Errorf("AllSeverities() returned %d severities, want 4", len(sevs))
	}
}

func TestParseSeverityValid(t *testing.T) {
	for _, raw := range []string{"critical", "major", "minor", "info"} {
		sev, err := model.ParseSeverity(raw)
		if err != nil {
			t.Errorf("ParseSeverity(%q) returned unexpected error: %v", raw, err)
		}
		if sev.String() != raw {
			t.Errorf("ParseSeverity(%q).String() = %q", raw, sev.String())
		}
	}
}

func TestParseSeverityInvalid(t *testing.T) {
	_, err := model.ParseSeverity("unknown")
	if err == nil {
		t.Error("ParseSeverity(\"unknown\") expected error, got nil")
	}
}
