package service

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/report"
	"errors"
	"strings"
)

type CheckResult struct {
	RecordID string
	Errors   []string
	Warnings []string
}

func CheckRecord(record domain.Record) CheckResult {
	result := CheckResult{RecordID: record.ID, Errors: []string{}, Warnings: []string{}}
	if err := record.Validate(); err != nil {
		result.Errors = append(result.Errors, err.Error())
	}
	if record.RequiresPhoto() && len(record.PhotoRefs) == 0 {
		result.Warnings = append(result.Warnings, "photo evidence missing")
	}
	if strings.TrimSpace(record.Description) == "" {
		result.Warnings = append(result.Warnings, "description missing")
	}
	return result
}
func (c CheckResult) Valid() bool { return len(c.Errors) == 0 }
func (c CheckResult) Messages() []string {
	return append(append([]string{}, c.Errors...), c.Warnings...)
}
func (s *Service) ValidateRecord(principal, id string) (CheckResult, error) {
	record, err := s.authorizedRecord(principal, "view_todo", id)
	if err != nil {
		return CheckResult{}, err
	}
	return CheckRecord(record), nil
}
func (s *Service) ValidateAll(principal string) ([]CheckResult, error) {
	records, err := s.Search(principal, domain.RecordFilter{IncludeArchived: true})
	if err != nil {
		return nil, err
	}
	out := make([]CheckResult, 0, len(records))
	for _, record := range records {
		out = append(out, CheckRecord(record))
	}
	return out, nil
}
func (s *Service) Findings(principal string) ([]report.Finding, error) {
	records, err := s.Search(principal, domain.RecordFilter{IncludeArchived: true})
	if err != nil {
		return nil, err
	}
	return report.ValidateRecords(records), nil
}
func (s *Service) RequireValid(record domain.Record) error {
	result := CheckRecord(record)
	if !result.Valid() {
		return errors.New(strings.Join(result.Errors, "; "))
	}
	return nil
}
