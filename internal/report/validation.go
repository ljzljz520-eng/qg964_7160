package report

import (
	"coldrisk.local/console/internal/domain"
	"strings"
)

type Finding struct {
	RecordID string
	Level    string
	Message  string
}

func ValidateRecords(records []domain.Record) []Finding {
	findings := []Finding{}
	for _, record := range records {
		if err := record.Validate(); err != nil {
			findings = append(findings, Finding{RecordID: record.ID, Level: "error", Message: err.Error()})
		}
		if record.RequiresPhoto() && len(record.PhotoRefs) == 0 {
			findings = append(findings, Finding{RecordID: record.ID, Level: "warning", Message: "high severity lacks photo evidence"})
		}
		if strings.TrimSpace(record.Description) == "" {
			findings = append(findings, Finding{RecordID: record.ID, Level: "warning", Message: "description is empty"})
		}
	}
	return findings
}
func Errors(findings []Finding) []Finding {
	out := []Finding{}
	for _, item := range findings {
		if item.Level == "error" {
			out = append(out, item)
		}
	}
	return out
}
func Warnings(findings []Finding) []Finding {
	out := []Finding{}
	for _, item := range findings {
		if item.Level == "warning" {
			out = append(out, item)
		}
	}
	return out
}
func FindingText(findings []Finding) string {
	lines := make([]string, 0, len(findings))
	for _, item := range findings {
		lines = append(lines, item.Level+":"+item.RecordID+":"+item.Message)
	}
	return strings.Join(lines, "\n")
}
