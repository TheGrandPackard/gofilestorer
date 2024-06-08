package gofilestorer

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

type jsonReader[V reader] struct {
	storer[V]
}

// Create a new reader that is backed by a JSON file
func NewJSONReader[V reader](fs afero.Fs, fileName string) (Reader[V], error) {
	s := &jsonReader[V]{
		storer: storer[V]{
			fs:       fs,
			fileName: fileName,
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

// read the file into the storer
func (s *jsonReader[V]) readFile() error {
	// Read file from disk
	dataBytes, err := afero.ReadFile(s.fs, s.fileName)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal JSON to struct
	data := []V{}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}
	s.data = data

	// Create map of data
	dataMap := map[uuid.UUID]*V{}
	for _, record := range data {
		dataMap[record.GetID()] = &record
	}
	s.dataMap = dataMap

	return nil
}
