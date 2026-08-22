package store

import (
	"errors"

	"coldrisk.local/console/internal/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) SaveRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	data, err := record.Marshal()
	if err != nil {
		return err
	}
	return s.Update(func(tx *bolt.Tx) error { return putValue(tx, recordsBucket, []byte(record.ID), data) })
}

func putValue(tx *bolt.Tx, bucket, key, value []byte) error {
	b := tx.Bucket(bucket)
	if b == nil {
		return errors.New("bucket missing")
	}
	return b.Put(key, value)
}

func (s *Store) GetRecord(id string) (domain.Record, error) {
	var r domain.Record
	if id == "" {
		return r, errors.New("record id is required")
	}
	err := s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(recordsBucket)
		raw := b.Get([]byte(id))
		if len(raw) == 0 {
			return errors.New("record not found")
		}
		return r.Unmarshal(raw)
	})
	return r, err
}
func (s *Store) ListRecords() ([]domain.Record, error) {
	out := []domain.Record{}
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(recordsBucket).ForEach(func(_, v []byte) error {
			var r domain.Record
			if err := r.Unmarshal(v); err != nil {
				return err
			}
			out = append(out, r)
			return nil
		})
	})
	domain.SortRecords(out)
	return out, err
}
func (s *Store) DeleteRecord(id string) error {
	if id == "" {
		return errors.New("record id is required")
	}
	return s.Update(func(tx *bolt.Tx) error { return tx.Bucket(recordsBucket).Delete([]byte(id)) })
}
func (s *Store) CountRecords() (int, error) { items, err := s.ListRecords(); return len(items), err }
func (s *Store) ReplaceRecord(record domain.Record) error {
	if err := record.Validate(); err != nil {
		return err
	}
	return s.SaveRecord(record)
}
