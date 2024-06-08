package gofilestorer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

type jsonWriter[V writer] struct {
	jsonReader[V]
}

// Create a new writer that is backed by a JSON file
func NewJSONWriter[V writer](fs afero.Fs, fileName string) (Writer[V], error) {
	s := &jsonWriter[V]{
		jsonReader: jsonReader[V]{
			storer: storer[V]{
				fs:       fs,
				fileName: fileName,
			},
		},
	}

	// Read file
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if err := s.readFile(); err != nil {
		return nil, err
	}

	return s, nil
}

// write the file from the storer
func (s *jsonWriter[V]) writeFile() error {
	// Marshal JSON to bytes
	dataBytes, err := json.Marshal(s.data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	// Write file to disk
	if err := afero.WriteFile(s.fs, s.fileName, dataBytes, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

// create a new record in the storer and write changes to file
func (s *jsonWriter[V]) Create(data V) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := uuid.New()
	data.SetID(id)
	data.SetCreatedAt(time.Now())
	s.data = append(s.data, data)
	s.dataMap[id] = &data

	return s.writeFile()
}

// update an existing record in the storer and write changes to file
func (s *jsonWriter[V]) Update(id uuid.UUID, data V) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.dataMap[id]
	if ok {
		data.SetUpdatedAt(time.Now())
		s.dataMap[data.GetID()] = &data
		return s.writeFile()
	}

	return ErrorDataNotExists
}

// delete an existing record in the storer and write changes to file
func (s *jsonWriter[V]) Delete(id uuid.UUID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.dataMap[id]
	if ok {
		delete(s.dataMap, id)
		for i, data := range s.data {
			if data.GetID() == id {
				s.data = append(s.data[:i], s.data[i+1:]...)
				return s.writeFile()
			}
		}
	}

	return ErrorDataNotExists
}
