package domain

import (
	"errors"
	"strings"
)

func (r *Record) Transition(next string) error {
	if r == nil {
		return errors.New("record is nil")
	}
	if !ValidStatus(next) {
		return errors.New("invalid target status")
	}
	if r.Status == next {
		return nil
	}
	allowed := map[string]map[string]bool{
		StatusOpen:      {StatusInReview: true, StatusArchived: false},
		StatusInReview:  {StatusConfirmed: true, StatusOpen: true},
		StatusConfirmed: {StatusPublished: true, StatusInReview: true},
		StatusPublished: {StatusArchived: true},
		StatusArchived:  {},
	}
	if !allowed[r.Status][next] {
		return errors.New("illegal status transition")
	}
	r.Status = next
	return nil
}

func (r Record) IsActive() bool       { return r.Status != StatusArchived }
func (r Record) NeedsAttention() bool { return r.Status == StatusOpen || r.Status == StatusInReview }
func (r Record) SeverityRank() int {
	switch strings.ToLower(r.Severity) {
	case SeverityCritical:
		return 4
	case SeverityHigh:
		return 3
	case SeverityMedium:
		return 2
	case SeverityLow:
		return 1
	default:
		return 0
	}
}
func (r Record) DisplayLabel() string { return r.StoreID + " / " + r.Title + " [" + r.Severity + "]" }
func (r Record) HasPhoto(id string) bool {
	for _, p := range r.PhotoRefs {
		if p == id {
			return true
		}
	}
	return false
}
func (r *Record) AddPhoto(id string) error {
	if r == nil || id == "" {
		return errors.New("photo id is required")
	}
	if !r.HasPhoto(id) {
		r.PhotoRefs = append(r.PhotoRefs, id)
	}
	return nil
}
func (r *Record) Assign(actor string) error {
	if r == nil || strings.TrimSpace(actor) == "" {
		return errors.New("assignee is required")
	}
	r.Assignee = actor
	return nil
}
func (w *Workflow) Advance(stage string) error {
	if w == nil || stage == "" {
		return errors.New("stage is required")
	}
	if w.Stage == stage {
		return nil
	}
	w.Stage = stage
	w.State = "active"
	return nil
}
func (w Workflow) IsDue(reference string) bool {
	return w.DueAt != "" && w.DueAt <= reference && w.State != "closed"
}
func (w *Workflow) Close() {
	if w != nil {
		w.State = "closed"
	}
}
func (e AuditEvent) Key() string { return e.StoreID + ":" + e.RecordID + ":" + e.ID }
