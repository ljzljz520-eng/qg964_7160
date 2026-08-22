package store

import (
	"coldrisk.local/console/internal/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) Search(filter domain.RecordFilter) ([]domain.Record, error) {
	items, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	return domain.FilterRecords(items, filter), nil
}
func (s *Store) RecordsByStore(storeID string) ([]domain.Record, error) {
	return s.Search(domain.RecordFilter{StoreID: storeID, IncludeArchived: true})
}
func (s *Store) ActiveRecords() ([]domain.Record, error) { return s.Search(domain.RecordFilter{}) }
func (s *Store) AuditFor(recordID string) ([]domain.AuditEvent, error) {
	items, err := s.ListAudits()
	return domain.FilterAuditEvents(items, recordID), err
}
func (s *Store) Snapshot() (map[string]int, error) {
	records, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	attachments, e := s.ListAttachments()
	if e != nil {
		return nil, e
	}
	workflows, e := s.ListWorkflows()
	if e != nil {
		return nil, e
	}
	audits, e := s.ListAudits()
	if e != nil {
		return nil, e
	}
	return map[string]int{"records": len(records), "attachments": len(attachments), "workflows": len(workflows), "audits": len(audits)}, nil
}
func (s *Store) Clear() error {
	return s.Update(func(tx *bolt.Tx) error {
		for _, name := range [][]byte{recordsBucket, attachmentsBucket, workflowsBucket, auditsBucket} {
			b := tx.Bucket(name)
			keys := [][]byte{}
			if err := b.ForEach(func(k, _ []byte) error { keys = append(keys, append([]byte(nil), k...)); return nil }); err != nil {
				return err
			}
			for _, k := range keys {
				if err := b.Delete(k); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
