package gofilestorer

import (
	"bytes"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/trimmer-io/go-csv"
)

type csvReader[V reader] struct {
	storer[V]
	separator rune
}

// Create a new reader that is backed by a CSV file
func NewCSVReader[V reader](fs afero.Fs, fileName string, separator rune) (Reader[V], error) {
	s := &csvReader[V]{
		storer: storer[V]{
			fs:       fs,
			fileName: fileName,
		},
		separator: separator,
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
func (s *csvReader[V]) readFile() error {
	// Read file from disk
	dataBytes, err := afero.ReadFile(s.fs, s.fileName)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal CSV to struct
	data := []V{}
	decoder := csv.NewDecoder(bytes.NewReader(dataBytes))
	decoder.Separator(s.separator)
	err = decoder.Decode(&data)
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
