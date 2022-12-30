package filestorer

import (
	"fmt"
	"time"

	"github.com/spf13/afero"
	"github.com/trimmer-io/go-csv"
)

type csvStorer[D data] struct {
	storer[D]
}

// Create a new Timecard storer that is backed by a JSON file
func NewCSVStorer[D data](fs afero.Fs, fileName string) (Storer[D], error) {
	s := &csvStorer[D]{
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
func (s *csvStorer[D]) readFile() error {
	// Read file from disk
	dataBytes, err := afero.ReadFile(s.fs, s.fileName)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal JSON to struct
	data := []D{}
	err = csv.Unmarshal(dataBytes, &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}

	s.data = data

	return nil
}

// write the file from the storer
func (s *csvStorer[D]) writeFile() error {
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
func (s *csvStorer[D]) Create(data D) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data.SetID(uint64(len(s.data) + 1))
	data.SetCreatedAt(time.Now())
	s.data = append(s.data, data)

	return s.writeFile()
}

// read all records from the storer
func (s *csvStorer[D]) Read() ([]D, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.data, nil
}

// update an existing record in the storer and write changes to file
func (s *csvStorer[D]) Update(data D) error {
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
func (s *csvStorer[D]) Delete(id uint64) error {
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
