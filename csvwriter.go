package gofilestorer

import (
	"fmt"
	"time"

	"github.com/spf13/afero"
	"github.com/trimmer-io/go-csv"
)

type csvWriter[D writer] struct {
	csvReader[D]
}

// Create a new writer that is backed by a CSV file
func NewCSVWriter[D writer](fs afero.Fs, fileName string, separator rune) (Writer[D], error) {
	s := &csvWriter[D]{
		csvReader: csvReader[D]{
			storer: storer[D]{
				fs:       fs,
				fileName: fileName,
			},
			separator: separator,
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
func (s *csvWriter[D]) writeFile() error {
	// Marshal JSON to bytes
	dataBytes, err := csv.Marshal(s.data)
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
func (s *csvWriter[D]) Create(data D) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if data.GetID() == 0 {
		data.SetID(uint64(len(s.data) + 1))
	}
	data.SetCreatedAt(time.Now())
	s.data = append(s.data, data)

	return s.writeFile()
}

// update an existing record in the storer and write changes to file
func (s *csvWriter[D]) Update(data D) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, record := range s.data {
		if record.GetID() == data.GetID() {
			s.data[i] = data
			return s.writeFile()
		}
	}

	return ErrorDataNotExists
}

// delete an existing record in the storer and write changes to file
func (s *csvWriter[D]) Delete(id uint64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, data := range s.data {
		if data.GetID() == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return s.writeFile()
		}
	}

	return ErrorDataNotExists
}

func (s *csvWriter[D]) Upsert(data D) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Attempt update
	for i, record := range s.data {
		if record.GetID() == data.GetID() {
			s.data[i] = data
			return s.writeFile()
		}
	}

	// Fall back to create
	if data.GetID() == 0 {
		data.SetID(uint64(len(s.data) + 1))
	}
	data.SetCreatedAt(time.Now())
	s.data = append(s.data, data)

	return s.writeFile()
}
