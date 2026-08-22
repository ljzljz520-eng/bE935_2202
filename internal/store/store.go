package store

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/bbolt"
)

var bucketNames = [][]byte{[]byte("cases"), []byte("files"), []byte("permissions"), []byte("versions"), []byte("audits"), []byte("shares")}

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0o600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, e := tx.CreateBucketIfNotExists(name); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func encode(value any) ([]byte, error)    { return json.Marshal(value) }
func decode(data []byte, value any) error { return json.Unmarshal(data, value) }

func (s *Store) Put(bucket, key string, value any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	data, err := encode(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s missing", bucket)
		}
		return b.Put([]byte(key), data)
	})
}

func (s *Store) Get(bucket, key string, value any) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return false, fmt.Errorf("store closed")
	}
	var found bool
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s missing", bucket)
		}
		raw := b.Get([]byte(key))
		if raw == nil {
			return nil
		}
		found = true
		return decode(raw, value)
	})
	return found, err
}

func (s *Store) Delete(bucket, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s missing", bucket)
		}
		return b.Delete([]byte(key))
	})
}

func (s *Store) List(bucket string) ([][]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, fmt.Errorf("store closed")
	}
	var rows [][]byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s missing", bucket)
		}
		return b.ForEach(func(_, v []byte) error {
			if v != nil {
				cp := append([]byte(nil), v...)
				rows = append(rows, cp)
			}
			return nil
		})
	})
	return rows, err
}

func (s *Store) Count(bucket string) (int, error) { rows, err := s.List(bucket); return len(rows), err }
