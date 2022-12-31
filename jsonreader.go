package filestorer

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
)

type jsonReader[D reader] struct {
	storer[D]
}

// Create a new reader that is backed by a JSON file
func NewJSONReader[D reader](fs afero.Fs, fileName string) (Reader[D], error) {
	s := &jsonReader[D]{
		storer: storer[D]{
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
func (s *jsonReader[D]) readFile() error {
	// Read file from disk
	dataBytes, err := afero.ReadFile(s.fs, s.fileName)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal JSON to struct
	data := []D{}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}

	s.data = data

	return nil
}

// read all records from the storer
func (s *jsonReader[D]) Read() ([]D, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.data, nil
}
