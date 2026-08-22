package service

import (
	"coldrisk.local/console/internal/domain"
	"errors"
	"strings"
)

type Collaboration struct {
	Record      domain.Record
	Attachments []domain.Attachment
	Workflow    domain.Workflow
	Audits      []domain.AuditEvent
}

func (s *Service) Collaborate(principal, id string) (Collaboration, error) {
	record, attachments, err := s.Select(principal, id)
	if err != nil {
		return Collaboration{}, err
	}
	workflow, e := s.Store.WorkflowFor(id)
	if e != nil {
		return Collaboration{}, e
	}
	audits, e := s.Store.AuditFor(id)
	if e != nil {
		return Collaboration{}, e
	}
	return Collaboration{Record: record, Attachments: attachments, Workflow: workflow, Audits: audits}, nil
}
func (s *Service) AddComment(principal, id, detail, at string) error {
	if strings.TrimSpace(detail) == "" {
		return errors.New("comment is required")
	}
	record, err := s.authorizedRecord(principal, "review_record", id)
	if err != nil {
		return err
	}
	return s.appendAudit(record, principal, "commented", detail, at)
}
func (s *Service) AssignRecord(principal, id, assignee, at string) (domain.Record, error) {
	record, err := s.authorizedRecord(principal, "review_record", id)
	if err != nil {
		return record, err
	}
	if err = record.Assign(assignee); err != nil {
		return record, err
	}
	record.UpdatedAt = at
	if err = s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendAudit(record, principal, "assigned", "assigned to "+assignee, at)
}
func (s *Service) PublishCollaboration(principal, id, at string) (Collaboration, error) {
	if _, err := s.Publish(principal, id, at); err != nil {
		return Collaboration{}, err
	}
	return s.Collaborate(principal, id)
}
