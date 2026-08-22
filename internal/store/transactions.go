package store

import (
	"coldrisk.local/console/internal/domain"
	"errors"
	bolt "go.etcd.io/bbolt"
)

type Bundle struct {
	Record     domain.Record
	Attachment domain.Attachment
	Workflow   domain.Workflow
	Audit      domain.AuditEvent
}

func (s *Store) SaveBundle(bundle Bundle) error {
	if err := domain.EnsureEntityValid(bundle.Record, bundle.Attachment, bundle.Workflow, bundle.Audit); err != nil {
		return err
	}
	return s.Update(func(tx *bolt.Tx) error {
		recordData, _ := bundle.Record.Marshal()
		attachmentData, _ := bundle.Attachment.Marshal()
		workflowData, _ := bundle.Workflow.Marshal()
		auditData, _ := bundle.Audit.Marshal()
		if err := tx.Bucket(recordsBucket).Put([]byte(bundle.Record.ID), recordData); err != nil {
			return err
		}
		if err := tx.Bucket(attachmentsBucket).Put([]byte(bundle.Attachment.ID), attachmentData); err != nil {
			return err
		}
		if err := tx.Bucket(workflowsBucket).Put([]byte(bundle.Workflow.ID), workflowData); err != nil {
			return err
		}
		return tx.Bucket(auditsBucket).Put([]byte(bundle.Audit.ID), auditData)
	})
}
func (s *Store) BundleFor(recordID string) (Bundle, error) {
	record, err := s.GetRecord(recordID)
	if err != nil {
		return Bundle{}, err
	}
	attachments, err := s.AttachmentsFor(recordID)
	if err != nil || len(attachments) == 0 {
		return Bundle{}, errors.New("bundle attachment missing")
	}
	workflow, err := s.WorkflowFor(recordID)
	if err != nil {
		return Bundle{}, err
	}
	audit, err := s.LatestAuditFor(recordID)
	if err != nil {
		return Bundle{}, err
	}
	return Bundle{Record: record, Attachment: attachments[0], Workflow: workflow, Audit: audit}, nil
}
func (s *Store) SaveMany(bundles []Bundle) error {
	for _, bundle := range bundles {
		if err := s.SaveBundle(bundle); err != nil {
			return err
		}
	}
	return nil
}
func (s *Store) DeleteBundle(recordID string) error {
	record, err := s.GetRecord(recordID)
	if err != nil {
		return err
	}
	if err = s.DeleteRecord(record.ID); err != nil {
		return err
	}
	attachments, _ := s.AttachmentsFor(recordID)
	for _, item := range attachments {
		_ = s.Update(func(tx *bolt.Tx) error { return tx.Bucket(attachmentsBucket).Delete([]byte(item.ID)) })
	}
	return nil
}
