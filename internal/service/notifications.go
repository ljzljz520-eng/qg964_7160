package service

import (
	"coldrisk.local/console/internal/domain"
	"sort"
)

type Notification struct {
	RecordID  string
	StoreID   string
	Recipient string
	Kind      string
	Message   string
}

func (s *Service) Notifications(principal string) ([]Notification, error) {
	items, err := s.Todo(principal, domain.RecordFilter{})
	if err != nil {
		return nil, err
	}
	out := make([]Notification, 0)
	for _, record := range items {
		recipient := record.Assignee
		if recipient == "" {
			recipient = principal
		}
		kind := "reminder"
		if record.IsCritical() {
			kind = "escalation"
		}
		out = append(out, Notification{RecordID: record.ID, StoreID: record.StoreID, Recipient: recipient, Kind: kind, Message: record.Title})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RecordID < out[j].RecordID })
	return out, nil
}
func (s *Service) NotificationsForStore(principal, storeID string) ([]Notification, error) {
	return s.NotificationsFiltered(principal, domain.RecordFilter{StoreID: storeID})
}
func (s *Service) NotificationsFiltered(principal string, filter domain.RecordFilter) ([]Notification, error) {
	items, err := s.Todo(principal, filter)
	if err != nil {
		return nil, err
	}
	out := make([]Notification, 0, len(items))
	for _, record := range items {
		kind := "reminder"
		if record.Severity == domain.SeverityCritical {
			kind = "escalation"
		}
		out = append(out, Notification{RecordID: record.ID, StoreID: record.StoreID, Recipient: record.Assignee, Kind: kind, Message: record.DisplayLabel()})
	}
	return out, nil
}
