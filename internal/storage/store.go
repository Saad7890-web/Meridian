package storage

import (
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	data    map[string]Entry
	wal     *WAL
	version int64
}

func NewStore(wal *WAL) *Store {
	return &Store{
		data: make(map[string]Entry),
		wal:  wal,
	}
}

func (s *Store) Put(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.version++

	entry := Entry{
		Key:     key,
		Value:   value,
		Version: s.version,
	}

	// 1. Write to WAL first
	if err := s.wal.Write(entry); err != nil {
		return err
	}

	// 2. Then update memory
	s.data[key] = entry

	return nil
}


func (s *Store) Get(key string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}


func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.version++

	entry := Entry{
		Key:     key,
		Value:   "",
		Version: s.version,
	}

	if err := s.wal.Write(entry); err != nil {
		return err
	}

	delete(s.data, key)
	return nil
}

func (s *Store) Recover() error {
	entries, err := s.wal.Load()
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Value == "" {
			delete(s.data, e.Key)
		} else {
			s.data[e.Key] = e
			if e.Version > s.version {
				s.version = e.Version
			}
		}
	}

	return nil
}