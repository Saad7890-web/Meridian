package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

type WAL struct {
	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
}

func NewWAL(path string) (*WAL, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (w *WAL) Write(entry Entry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = w.writer.Write(append(data, '\n'))
	if err != nil {
		return err
	}

	return w.writer.Flush()
}


func (w *WAL) Load() ([]Entry, error) {
	file, err := os.Open(w.file.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var e Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}

	return entries, scanner.Err()
}