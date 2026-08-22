package store

import (
	"coldrisk.local/console/internal/domain"
	"errors"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) SaveAttachment(item domain.Attachment) error {
	if err := item.Validate(); err != nil {
		return err
	}
	data, err := item.Marshal()
	if err != nil {
		return err
	}
	return s.Update(func(tx *bolt.Tx) error { return putValue(tx, attachmentsBucket, []byte(item.ID), data) })
}
func (s *Store) ListAttachments() ([]domain.Attachment, error) {
	out := []domain.Attachment{}
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(attachmentsBucket).ForEach(func(_, v []byte) error {
			var a domain.Attachment
			if err := a.Unmarshal(v); err != nil {
				return err
			}
			out = append(out, a)
			return nil
		})
	})
	domain.SortAttachments(out)
	return out, err
}
func (s *Store) AttachmentsFor(recordID string) ([]domain.Attachment, error) {
	items, err := s.ListAttachments()
	return domain.FilterAttachments(items, recordID), err
}
func (s *Store) SaveWorkflow(item domain.Workflow) error {
	if err := item.Validate(); err != nil {
		return err
	}
	data, err := item.Marshal()
	if err != nil {
		return err
	}
	return s.Update(func(tx *bolt.Tx) error { return putValue(tx, workflowsBucket, []byte(item.ID), data) })
}
func (s *Store) ListWorkflows() ([]domain.Workflow, error) {
	out := []domain.Workflow{}
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(workflowsBucket).ForEach(func(_, v []byte) error {
			var w domain.Workflow
			if err := w.Unmarshal(v); err != nil {
				return err
			}
			out = append(out, w)
			return nil
		})
	})
	domain.SortWorkflows(out)
	return out, err
}
func (s *Store) WorkflowFor(recordID string) (domain.Workflow, error) {
	items, err := s.ListWorkflows()
	for _, w := range items {
		if w.RecordID == recordID {
			return w, err
		}
	}
	return domain.Workflow{}, errors.New("workflow not found")
}
func (s *Store) SaveAudit(item domain.AuditEvent) error {
	if err := item.Validate(); err != nil {
		return err
	}
	data, err := item.Marshal()
	if err != nil {
		return err
	}
	return s.Update(func(tx *bolt.Tx) error { return putValue(tx, auditsBucket, []byte(item.ID), data) })
}
func (s *Store) ListAudits() ([]domain.AuditEvent, error) {
	out := []domain.AuditEvent{}
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(auditsBucket).ForEach(func(_, v []byte) error {
			var e domain.AuditEvent
			if err := e.Unmarshal(v); err != nil {
				return err
			}
			out = append(out, e)
			return nil
		})
	})
	domain.SortAuditEvents(out)
	return out, err
}
