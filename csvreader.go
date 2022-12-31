package filestorer

import (
	"bytes"
	"fmt"

	"github.com/spf13/afero"
	"github.com/trimmer-io/go-csv"
)

type csvReader[D reader] struct {
	storer[D]
	separator rune
}

// Create a new reader that is backed by a CSV file
func NewCSVReader[D reader](fs afero.Fs, fileName string, separator rune) (Reader[D], error) {
	s := &csvReader[D]{
		storer: storer[D]{
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
func (s *csvReader[D]) readFile() error {
	// Read file from disk
	dataBytes, err := afero.ReadFile(s.fs, s.fileName)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal JSON to struct
	data := []D{}
	decoder := csv.NewDecoder(bytes.NewReader(dataBytes))
	decoder.Separator(s.separator)
	err = decoder.Decode(&data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}

	s.data = data

	return nil
}

// read all records from the storer
func (s *csvReader[D]) Read() ([]D, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.data, nil
}
