package domain

import "strings"

type RecordFilter struct {
	StoreID         string
	Status          string
	Severity        string
	Query           string
	IncludeArchived bool
}

func (f RecordFilter) Match(r Record) bool {
	if f.StoreID != "" && r.StoreID != f.StoreID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Severity != "" && r.Severity != f.Severity {
		return false
	}
	if !f.IncludeArchived && r.Status == StatusArchived {
		return false
	}
	if f.Query != "" {
		q := strings.ToLower(f.Query)
		if !strings.Contains(strings.ToLower(r.Title), q) && !strings.Contains(strings.ToLower(r.Description), q) {
			return false
		}
	}
	return true
}
func FilterRecords(items []Record, filter RecordFilter) []Record {
	out := make([]Record, 0, len(items))
	for _, r := range items {
		if filter.Match(r) {
			out = append(out, CloneRecord(r))
		}
	}
	SortRecords(out)
	return out
}
func FilterAttachments(items []Attachment, recordID string) []Attachment {
	out := make([]Attachment, 0)
	for _, a := range items {
		if recordID == "" || a.RecordID == recordID {
			out = append(out, a)
		}
	}
	SortAttachments(out)
	return out
}
func FilterWorkflows(items []Workflow, storeID string) []Workflow {
	out := make([]Workflow, 0)
	for _, w := range items {
		if storeID == "" || w.StoreID == storeID {
			out = append(out, w)
		}
	}
	SortWorkflows(out)
	return out
}
func FilterAuditEvents(items []AuditEvent, recordID string) []AuditEvent {
	out := make([]AuditEvent, 0)
	for _, e := range items {
		if recordID == "" || e.RecordID == recordID {
			out = append(out, e)
		}
	}
	SortAuditEvents(out)
	return out
}
