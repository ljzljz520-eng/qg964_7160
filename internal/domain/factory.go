package domain

import "errors"

func NewRecord(id, store, title, severity, at string) Record {
	return Record{ID: id, StoreID: store, Title: title, Status: StatusOpen, Severity: severity, CreatedAt: at, UpdatedAt: at}
}
func NewAttachment(id, recordID, store, kind, uri, actor, at string) Attachment {
	return Attachment{ID: id, RecordID: recordID, StoreID: store, Kind: kind, URI: uri, UploadedBy: actor, UploadedAt: at}
}
func NewWorkflow(id, recordID, store, owner, due string) Workflow {
	return Workflow{ID: id, RecordID: recordID, StoreID: store, Owner: owner, DueAt: due, Stage: "intake", State: "active"}
}
func NewAuditEvent(id, recordID, actor, action, store, detail, at string) AuditEvent {
	return AuditEvent{ID: id, RecordID: recordID, Actor: actor, Action: action, StoreID: store, Detail: detail, At: at}
}
func CloneRecord(r Record) Record { r.PhotoRefs = append([]string(nil), r.PhotoRefs...); return r }
func CloneAttachments(items []Attachment) []Attachment {
	out := make([]Attachment, len(items))
	copy(out, items)
	return out
}
func CloneWorkflows(items []Workflow) []Workflow {
	out := make([]Workflow, len(items))
	copy(out, items)
	return out
}
func CloneAuditEvents(items []AuditEvent) []AuditEvent {
	out := make([]AuditEvent, len(items))
	copy(out, items)
	return out
}
func EnsureEntityValid(record Record, attachment Attachment, workflow Workflow, event AuditEvent) error {
	if err := record.Validate(); err != nil {
		return err
	}
	if err := attachment.Validate(); err != nil {
		return err
	}
	if err := workflow.Validate(); err != nil {
		return err
	}
	return event.Validate()
}
func RequireText(values ...string) error {
	for _, value := range values {
		if value == "" {
			return errors.New("required text missing")
		}
	}
	return nil
}
