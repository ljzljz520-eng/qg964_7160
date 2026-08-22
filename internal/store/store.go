package store

import (
	"errors"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

var (
	recordsBucket     = []byte("records")
	attachmentsBucket = []byte("attachments")
	workflowsBucket   = []byte("workflows")
	auditsBucket      = []byte("audits")
)

type Store struct {
	db   *bolt.DB
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bolt.Open(filepath.Clean(path), 0600, &bolt.Options{NoSync: true})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range [][]byte{recordsBucket, attachmentsBucket, workflowsBucket, auditsBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
func (s *Store) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}
func (s *Store) healthy() error {
	if s == nil || s.db == nil {
		return errors.New("store is closed")
	}
	return nil
}
func (s *Store) View(fn func(*bolt.Tx) error) error {
	if err := s.healthy(); err != nil {
		return err
	}
	return s.db.View(fn)
}
func (s *Store) Update(fn func(*bolt.Tx) error) error {
	if err := s.healthy(); err != nil {
		return err
	}
	return s.db.Update(fn)
}
