package store

import (
	"coldrisk.local/console/internal/domain"
	"errors"
	"sort"
)

type Health struct {
	Records     int
	Attachments int
	Workflows   int
	Audits      int
}

func (s *Store) Health() (Health, error) {
	snapshot, err := s.Snapshot()
	if err != nil {
		return Health{}, err
	}
	return Health{Records: snapshot["records"], Attachments: snapshot["attachments"], Workflows: snapshot["workflows"], Audits: snapshot["audits"]}, nil
}
func (s *Store) ValidateAll() []error {
	errorsOut := []error{}
	records, err := s.ListRecords()
	if err != nil {
		return []error{err}
	}
	for _, record := range records {
		if err := record.Validate(); err != nil {
			errorsOut = append(errorsOut, err)
		}
	}
	attachments, err := s.ListAttachments()
	if err != nil {
		return append(errorsOut, err)
	}
	for _, item := range attachments {
		if err := item.Validate(); err != nil {
			errorsOut = append(errorsOut, err)
		}
	}
	workflows, err := s.ListWorkflows()
	if err != nil {
		return append(errorsOut, err)
	}
	for _, item := range workflows {
		if err := item.Validate(); err != nil {
			errorsOut = append(errorsOut, err)
		}
	}
	audits, err := s.ListAudits()
	if err != nil {
		return append(errorsOut, err)
	}
	for _, item := range audits {
		if err := item.Validate(); err != nil {
			errorsOut = append(errorsOut, err)
		}
	}
	return errorsOut
}
func (s *Store) RecordIDs() ([]string, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(records))
	for _, record := range records {
		out = append(out, record.ID)
	}
	sort.Strings(out)
	return out, nil
}
func (s *Store) Stores() ([]string, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	set := map[string]bool{}
	for _, record := range records {
		set[record.StoreID] = true
	}
	out := make([]string, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	sort.Strings(out)
	return out, nil
}
func (s *Store) PurgeArchived() (int, error) {
	records, err := s.ListRecords()
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, record := range records {
		if record.Status == domain.StatusArchived {
			if err := s.DeleteRecord(record.ID); err != nil {
				return removed, err
			}
			removed++
		}
	}
	return removed, nil
}
func (s *Store) ReplaceAll(records []domain.Record) error {
	for _, record := range records {
		if err := s.SaveRecord(record); err != nil {
			return err
		}
	}
	return nil
}
func (s *Store) CopyRecord(from, to string) error {
	record, err := s.GetRecord(from)
	if err != nil {
		return err
	}
	if to == "" {
		return errors.New("destination required")
	}
	record.ID = to
	return s.SaveRecord(record)
}
func (s *Store) RecordsNeedingReview() ([]domain.Record, error) {
	items, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	out := []domain.Record{}
	for _, record := range items {
		if record.NeedsAttention() {
			out = append(out, record)
		}
	}
	domain.SortRecords(out)
	return out, nil
}
func (s *Store) RecordsBySeverity(severity string) ([]domain.Record, error) {
	return s.Search(domain.RecordFilter{Severity: severity, IncludeArchived: true})
}
func (s *Store) LatestAuditFor(recordID string) (domain.AuditEvent, error) {
	items, err := s.AuditFor(recordID)
	if err != nil {
		return domain.AuditEvent{}, err
	}
	if len(items) == 0 {
		return domain.AuditEvent{}, errors.New("audit not found")
	}
	return items[len(items)-1], nil
}
