package service

import (
	"errors"
	"fmt"
	"strings"

	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/policy"
	"coldrisk.local/console/internal/report"
	"coldrisk.local/console/internal/store"
)

type Service struct {
	Store  *store.Store
	Policy *policy.Policy
	Report report.Builder
}

func New(st *store.Store, pol *policy.Policy) *Service {
	return &Service{Store: st, Policy: pol, Report: report.NewBuilder()}
}
func (s *Service) ready() error {
	if s == nil || s.Store == nil {
		return errors.New("service store is required")
	}
	if s.Policy == nil {
		return errors.New("service policy is required")
	}
	return nil
}
func (s *Service) CreateRecord(principal string, record domain.Record, attachment domain.Attachment, workflow domain.Workflow, at string) error {
	if err := s.ready(); err != nil {
		return err
	}
	if !s.Policy.CheckRecord(principal, policy.ActionCreate, record) {
		return errors.New("create forbidden")
	}
	if err := record.Validate(); err != nil {
		return err
	}
	if attachment.RecordID == "" {
		attachment = domain.NewAttachment("att-"+record.ID, record.ID, record.StoreID, "photo", "pending://"+record.ID, principal, at)
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return err
	}
	if err := s.Store.SaveAttachment(attachment); err != nil {
		return err
	}
	if workflow.ID == "" {
		workflow = domain.NewWorkflow("wf-"+record.ID, record.ID, record.StoreID, principal, "")
	}
	if err := s.Store.SaveWorkflow(workflow); err != nil {
		return err
	}
	return s.appendAudit(record, principal, "created", "record created", at)
}
func (s *Service) Review(principal, id, at string) (domain.Record, error) {
	record, err := s.authorizedRecord(principal, policy.ActionReview, id)
	if err != nil {
		return record, err
	}
	if err = record.Transition(domain.StatusInReview); err != nil {
		return record, err
	}
	record.UpdatedAt = at
	if err = s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	err = s.appendAudit(record, principal, "reviewed", "review started", at)
	return record, err
}
func (s *Service) Confirm(principal, id, at string) (domain.Record, error) {
	record, err := s.authorizedRecord(principal, policy.ActionReview, id)
	if err != nil {
		return record, err
	}
	if err = record.Transition(domain.StatusConfirmed); err != nil {
		return record, err
	}
	record.UpdatedAt = at
	if err = s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	err = s.appendAudit(record, principal, "confirmed", "remediation confirmed", at)
	return record, err
}
func (s *Service) Archive(principal, id, at string) (domain.Record, error) {
	record, err := s.authorizedRecord(principal, policy.ActionArchive, id)
	if err != nil {
		return record, err
	}
	if err = record.Transition(domain.StatusArchived); err != nil {
		return record, err
	}
	record.UpdatedAt = at
	if err = s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	workflow, e := s.Store.WorkflowFor(id)
	if e == nil {
		workflow.Close()
		_ = s.Store.SaveWorkflow(workflow)
	}
	err = s.appendAudit(record, principal, "archived", "record archived", at)
	return record, err
}
func (s *Service) Publish(principal, id, at string) (domain.Record, error) {
	record, err := s.authorizedRecord(principal, policy.ActionPublish, id)
	if err != nil {
		return record, err
	}
	if err = record.Transition(domain.StatusPublished); err != nil {
		return record, err
	}
	record.UpdatedAt = at
	if err = s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	err = s.appendAudit(record, principal, "published", "record published", at)
	return record, err
}
func (s *Service) authorizedRecord(principal, action, id string) (domain.Record, error) {
	if err := s.ready(); err != nil {
		return domain.Record{}, err
	}
	record, err := s.Store.GetRecord(id)
	if err != nil {
		return record, err
	}
	if !s.Policy.CheckRecord(principal, action, record) {
		return record, errors.New("record action forbidden")
	}
	return record, nil
}
func (s *Service) appendAudit(record domain.Record, actor, action, detail, at string) error {
	event := domain.NewAuditEvent(fmt.Sprintf("audit-%s-%s", record.ID, action), record.ID, actor, action, record.StoreID, detail, at)
	return s.Store.SaveAudit(event)
}
func (s *Service) AddAttachment(principal string, item domain.Attachment, at string) error {
	if err := s.ready(); err != nil {
		return err
	}
	record, err := s.Store.GetRecord(item.RecordID)
	if err != nil {
		return err
	}
	if !s.Policy.CheckRecord(principal, policy.ActionReview, record) {
		return errors.New("attachment forbidden")
	}
	if item.UploadedAt == "" {
		item.UploadedAt = at
	}
	return s.Store.SaveAttachment(item)
}
func (s *Service) UpdateDescription(principal, id, description, at string) (domain.Record, error) {
	record, err := s.authorizedRecord(principal, policy.ActionReview, id)
	if err != nil {
		return record, err
	}
	if strings.TrimSpace(description) == "" {
		return record, errors.New("description is required")
	}
	record.Description = description
	record.UpdatedAt = at
	if err = s.Store.SaveRecord(record); err != nil {
		return record, err
	}
	return record, s.appendAudit(record, principal, "updated", "description updated", at)
}
