package service

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/policy"
	"errors"
	"sort"
	"strings"
)

type Queue struct {
	Principal  string
	Records    []domain.Record
	Total      int
	Critical   int
	Unassigned int
}

func (s *Service) BuildQueue(principal string, filter domain.RecordFilter) (Queue, error) {
	items, err := s.Todo(principal, filter)
	if err != nil {
		return Queue{}, err
	}
	queue := Queue{Principal: principal, Records: items, Total: len(items)}
	for _, record := range items {
		if record.IsCritical() {
			queue.Critical++
		}
		if record.Assignee == "" {
			queue.Unassigned++
		}
	}
	return queue, nil
}
func (s *Service) Triage(principal, id, assignee, description, at string) (domain.Record, error) {
	record, err := s.authorizedRecord(principal, policy.ActionReview, id)
	if err != nil {
		return record, err
	}
	if err = record.Assign(assignee); err != nil {
		return record, err
	}
	record.Description = strings.TrimSpace(description)
	if record.Description == "" {
		return record, errors.New("triage description required")
	}
	if err = record.AddPhoto("triage-" + record.ID); err != nil {
		return record, err
	}
	record.UpdatedAt = at
	if err = s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendAudit(record, principal, "triaged", "triage completed", at)
}
func (s *Service) ConfirmEvidence(principal, id, photoID, at string) (domain.Record, error) {
	record, err := s.authorizedRecord(principal, policy.ActionReview, id)
	if err != nil {
		return record, err
	}
	if err = record.AddPhoto(photoID); err != nil {
		return record, err
	}
	record.UpdatedAt = at
	if err = s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendAudit(record, principal, "evidence_added", "photo evidence added", at)
}
func (s *Service) Reopen(principal, id, at string) (domain.Record, error) {
	record, err := s.authorizedRecord(principal, policy.ActionReview, id)
	if err != nil {
		return record, err
	}
	if record.Status == domain.StatusArchived {
		return record, errors.New("archived records cannot reopen")
	}
	if record.Status == domain.StatusConfirmed {
		record.Status = domain.StatusInReview
	}
	if record.Status == domain.StatusPublished {
		record.Status = domain.StatusConfirmed
	}
	record.UpdatedAt = at
	if err = s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendAudit(record, principal, "reopened", "record returned for remediation", at)
}
func (s *Service) BulkAssign(principal, assignee string, ids []string, at string) ([]domain.Record, error) {
	if strings.TrimSpace(assignee) == "" {
		return nil, errors.New("assignee required")
	}
	out := make([]domain.Record, 0, len(ids))
	for _, id := range ids {
		record, err := s.AssignRecord(principal, id, assignee, at)
		if err != nil {
			return out, err
		}
		out = append(out, record)
	}
	return out, nil
}
func (s *Service) BulkPublish(principal string, ids []string, at string) ([]domain.Record, error) {
	out := make([]domain.Record, 0, len(ids))
	for _, id := range ids {
		record, err := s.Publish(principal, id, at)
		if err != nil {
			return out, err
		}
		out = append(out, record)
	}
	return out, nil
}
func (s *Service) ActiveBySeverity(principal, severity string) ([]domain.Record, error) {
	items, err := s.Search(principal, domain.RecordFilter{Severity: severity})
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].SeverityRank() > items[j].SeverityRank() })
	return items, nil
}
func (s *Service) AuthorizedStores(principal string) []string {
	if s == nil || s.Policy == nil {
		return nil
	}
	return s.Policy.FilterStores(principal)
}
func (s *Service) Can(principal, action, storeID string) bool {
	return s != nil && s.Policy != nil && s.Policy.Can(principal, action, storeID)
}
func (s *Service) ExplainAccess(principal, action, storeID string) string {
	if s == nil || s.Policy == nil {
		return "unavailable"
	}
	return s.Policy.Explain(principal, action, storeID)
}
