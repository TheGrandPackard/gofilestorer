package gofilestorer

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

type storer[V any] struct {
	fs       afero.Fs
	fileName string
	mutex    sync.RWMutex
	data     []V
	dataMap  map[uuid.UUID]*V
}

type Reader[V reader] interface {
	readFile() error

	ReadAll() ([]V, error)
	ReadOne(uuid.UUID) (*V, error)
}

type reader interface {
	GetID() uuid.UUID
}

// read all records from the storer
func (s *storer[V]) ReadAll() ([]V, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data, nil
}

// read a record from the storer
func (s *storer[V]) ReadOne(id uuid.UUID) (*V, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.dataMap[id]
	if ok {

		return s.dataMap[id], nil
	}

	return nil, ErrorDataNotExists
}

type Writer[V writer] interface {
	Reader[V]
	writeFile() error

	Create(V) error
	Update(uuid.UUID, V) error
	Delete(uuid.UUID) error
}

type writer interface {
	GetID() uuid.UUID
	SetID(uuid.UUID)
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}
