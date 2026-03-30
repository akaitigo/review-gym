package model

import "fmt"

// Severity represents the severity level of a reference review point.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityMajor    Severity = "major"
	SeverityMinor    Severity = "minor"
	SeverityInfo     Severity = "info"
)

// AllSeverities returns all valid severities ordered from most to least severe.
func AllSeverities() []Severity {
	return []Severity{
		SeverityCritical,
		SeverityMajor,
		SeverityMinor,
		SeverityInfo,
	}
}

// Valid reports whether s is a recognized severity.
func (s Severity) Valid() bool {
	switch s {
	case SeverityCritical, SeverityMajor, SeverityMinor, SeverityInfo:
		return true
	default:
		return false
	}
}

// String returns the string representation of the severity.
func (s Severity) String() string {
	return string(s)
}

// ParseSeverity converts a raw string to a Severity, returning an error for unknown values.
func ParseSeverity(raw string) (Severity, error) {
	s := Severity(raw)
	if !s.Valid() {
		return "", fmt.Errorf("unknown severity: %q", raw)
	}
	return s, nil
}
