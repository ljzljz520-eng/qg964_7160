package service

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/report"
)

func (s *Service) BuildReport(principal string, filter domain.RecordFilter) (report.Summary, error) {
	items, err := s.Search(principal, filter)
	if err != nil {
		return report.Summary{}, err
	}
	audits, e := s.Store.ListAudits()
	if e != nil {
		return report.Summary{}, e
	}
	return s.Report.Build(items, audits), nil
}
func (s *Service) ExportReport(principal string, filter domain.RecordFilter) (string, error) {
	summary, err := s.BuildReport(principal, filter)
	if err != nil {
		return "", err
	}
	return s.Report.Format(summary), nil
}
func (s *Service) Health() bool { return s != nil && s.Store != nil && s.Policy != nil }
