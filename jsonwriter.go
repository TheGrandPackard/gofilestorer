package gofilestorer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/afero"
)

type jsonWriter[K comparable, V writer[K]] struct {
	jsonReader[K, V]
}

// Create a new writer that is backed by a JSON file
func NewJSONWriter[K comparable, V writer[K]](fs afero.Fs, fileName string, newIDFunc func(data []V) K) (Writer[K, V], error) {
	s := &jsonWriter[K, V]{
		jsonReader: jsonReader[K, V]{
			storer: storer[K, V]{
				fs:        fs,
				fileName:  fileName,
				newIDFunc: newIDFunc,
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
func (s *jsonWriter[K, V]) writeFile() error {
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
func (s *jsonWriter[K, V]) Create(data V) (V, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := s.newIDFunc(s.data)
	data.SetID(id)
	data.SetCreatedAt(time.Now())
	s.data = append(s.data, data)
	s.dataMap[id] = data

	return data, s.writeFile()
}

// update an existing record in the storer and write changes to file
func (s *jsonWriter[K, V]) Update(id K, data V) (V, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, ok := s.dataMap[id]
	if ok {
		data.SetUpdatedAt(time.Now())
		s.dataMap[data.GetID()] = data
		for _, d := range s.data {
			if d.GetID() == id {
				return data, s.writeFile()
			}
		}
	}

	return *new(V), ErrorDataNotExists
}

// delete an existing record in the storer and write changes to file
func (s *jsonWriter[K, V]) Delete(id K) error {
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
