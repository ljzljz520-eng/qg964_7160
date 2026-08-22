package service

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/policy"
	"errors"
)

func (s *Service) Todo(principal string, filter domain.RecordFilter) ([]domain.Record, error) {
	if err := s.ready(); err != nil {
		return nil, err
	}
	actions := s.Policy.AuthorizedActionSet(principal)
	if !actions[policy.ActionViewTodo] && !actions["*"] {
		return nil, errors.New("todo forbidden")
	}
	// The todo list must only surface remediation records for stores the
	// principal is authorized to see. A regional supervisor granted a single
	// store would otherwise see every store's remediation photos and assignees.
	if filter.StoreID != "" && !s.Policy.HasStore(principal, filter.StoreID) {
		return nil, errors.New("todo forbidden")
	}
	filter.IncludeArchived = false
	items, err := s.Store.Search(filter)
	if err != nil {
		return nil, err
	}
	return s.Policy.StoresForRecord(principal, items), nil
}
func (s *Service) TodoForStore(principal, storeID string) ([]domain.Record, error) {
	return s.Todo(principal, domain.RecordFilter{StoreID: storeID})
}
func (s *Service) TodoCount(principal string) (int, error) {
	items, err := s.Todo(principal, domain.RecordFilter{})
	return len(items), err
}
func (s *Service) PrioritizedTodo(principal string) ([]domain.Record, error) {
	items, err := s.Todo(principal, domain.RecordFilter{})
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].SeverityRank() > items[i].SeverityRank() {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return items, nil
}
