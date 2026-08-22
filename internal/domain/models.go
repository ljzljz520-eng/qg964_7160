package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

const (
	StatusOpen       = "open"
	StatusInReview   = "in_review"
	StatusConfirmed  = "confirmed"
	StatusPublished  = "published"
	StatusArchived   = "archived"
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)

type Record struct {
	ID          string   `json:"id"`
	StoreID     string   `json:"store_id"`
	Title       string   `json:"title"`
	Status      string   `json:"status"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	PhotoRefs   []string `json:"photo_refs"`
	Assignee    string   `json:"assignee"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type Attachment struct {
	ID         string `json:"id"`
	RecordID   string `json:"record_id"`
	StoreID    string `json:"store_id"`
	Kind       string `json:"kind"`
	URI        string `json:"uri"`
	Hash       string `json:"hash"`
	UploadedBy string `json:"uploaded_by"`
	UploadedAt string `json:"uploaded_at"`
}

type Workflow struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Stage    string `json:"stage"`
	Owner    string `json:"owner"`
	StoreID  string `json:"store_id"`
	DueAt    string `json:"due_at"`
	State    string `json:"state"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Actor    string `json:"actor"`
	Action   string `json:"action"`
	StoreID  string `json:"store_id"`
	Detail   string `json:"detail"`
	At       string `json:"at"`
}

type Permission struct {
	Principal string            `json:"principal"`
	Roles     []string          `json:"roles"`
	Stores    map[string]string `json:"stores"`
	Actions   map[string]bool   `json:"actions"`
}

func (r Record) Marshal() ([]byte, error)         { return json.Marshal(r) }
func (r *Record) Unmarshal(data []byte) error     { return json.Unmarshal(data, r) }
func (a Attachment) Marshal() ([]byte, error)     { return json.Marshal(a) }
func (a *Attachment) Unmarshal(data []byte) error { return json.Unmarshal(data, a) }
func (w Workflow) Marshal() ([]byte, error)       { return json.Marshal(w) }
func (w *Workflow) Unmarshal(data []byte) error   { return json.Unmarshal(data, w) }
func (e AuditEvent) Marshal() ([]byte, error)     { return json.Marshal(e) }
func (e *AuditEvent) Unmarshal(data []byte) error { return json.Unmarshal(data, e) }

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id is required")
	}
	if strings.TrimSpace(r.StoreID) == "" {
		return errors.New("record store is required")
	}
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("record title is required")
	}
	if !ValidStatus(r.Status) {
		return fmt.Errorf("unsupported status %q", r.Status)
	}
	if !ValidSeverity(r.Severity) {
		return fmt.Errorf("unsupported severity %q", r.Severity)
	}
	if strings.TrimSpace(r.CreatedAt) == "" || strings.TrimSpace(r.UpdatedAt) == "" {
		return errors.New("record timestamps are required")
	}
	return nil
}

func (a Attachment) Validate() error {
	if a.ID == "" || a.RecordID == "" || a.StoreID == "" {
		return errors.New("attachment identity is required")
	}
	if a.Kind == "" || a.URI == "" {
		return errors.New("attachment kind and uri are required")
	}
	if a.UploadedAt == "" {
		return errors.New("attachment timestamp is required")
	}
	return nil
}

func (w Workflow) Validate() error {
	if w.ID == "" || w.RecordID == "" || w.StoreID == "" {
		return errors.New("workflow identity is required")
	}
	if w.Stage == "" || w.State == "" {
		return errors.New("workflow stage and state are required")
	}
	return nil
}

func (e AuditEvent) Validate() error {
	if e.ID == "" || e.RecordID == "" || e.Actor == "" || e.StoreID == "" {
		return errors.New("audit identity is required")
	}
	if e.Action == "" || e.At == "" {
		return errors.New("audit action and time are required")
	}
	return nil
}

func ValidStatus(status string) bool {
	switch status {
	case StatusOpen, StatusInReview, StatusConfirmed, StatusPublished, StatusArchived:
		return true
	default:
		return false
	}
}

func ValidSeverity(severity string) bool {
	switch severity {
	case SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical:
		return true
	default:
		return false
	}
}

func SortRecords(records []Record) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].StoreID == records[j].StoreID {
			return records[i].ID < records[j].ID
		}
		return records[i].StoreID < records[j].StoreID
	})
}

func SortAttachments(items []Attachment) {
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
}
func SortWorkflows(items []Workflow) {
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
}
func SortAuditEvents(items []AuditEvent) {
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
}
