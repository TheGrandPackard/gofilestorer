package gofilestorer

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
)

type jsonReader[K comparable, V reader[K]] struct {
	storer[K, V]
}

// Create a new reader that is backed by a JSON file
func NewJSONReader[K comparable, V reader[K]](fs afero.Fs, fileName string) (Reader[K, V], error) {
	s := &jsonReader[K, V]{
		storer: storer[K, V]{
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
func (s *jsonReader[K, V]) readFile() error {
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
	dataMap := map[K]*V{}
	for _, record := range data {
		dataMap[record.GetID()] = &record
	}
	s.dataMap = dataMap

	return nil
}
