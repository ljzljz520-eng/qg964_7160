package service

import (
	"coldrisk.local/console/internal/domain"
	"errors"
	"strings"
)

func (s *Service) Search(principal string, filter domain.RecordFilter) ([]domain.Record, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	if !s.Policy.Can(principal, "view_todo", filter.StoreID) && !s.Policy.AuthorizedActionSet(principal)["*"] {
		return nil, errors.New("search forbidden")
	}
	return s.Store.Search(filter)
}
func (s *Service) Select(principal, id string) (domain.Record, []domain.Attachment, error) {
	record, err := s.authorizedRecord(principal, "view_todo", id)
	if err != nil {
		return record, nil, err
	}
	attachments, e := s.Store.AttachmentsFor(id)
	return record, attachments, e
}
func (s *Service) MatchText(principal, query string) ([]domain.Record, error) {
	query = strings.TrimSpace(query)
	return s.Search(principal, domain.RecordFilter{Query: query})
}
func (s *Service) StoreSummary(principal, storeID string) (map[string]int, error) {
	items, err := s.Search(principal, domain.RecordFilter{StoreID: storeID, IncludeArchived: true})
	if err != nil {
		return nil, err
	}
	out := map[string]int{}
	for _, item := range items {
		out[item.Status]++
	}
	return out, nil
}
func (s *Service) RefreshWorkflow(principal, id, stage string) error {
	record, err := s.authorizedRecord(principal, "review_record", id)
	if err != nil {
		return err
	}
	workflow, err := s.Store.WorkflowFor(record.ID)
	if err != nil {
		return err
	}
	if err = workflow.Advance(stage); err != nil {
		return err
	}
	return s.Store.SaveWorkflow(workflow)
}
