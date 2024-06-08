package gofilestorer

import (
	"sync"
	"time"

	"github.com/spf13/afero"
)

type storer[K comparable, V any] struct {
	fs        afero.Fs
	fileName  string
	mutex     sync.RWMutex
	data      []V
	dataMap   map[K]*V
	newIDFunc func(data []V) K
}

type Reader[K comparable, V reader[K]] interface {
	readFile() error

	ReadAll() ([]V, error)
	ReadOne(K) (*V, error)
}

type reader[K comparable] interface {
	GetID() K
}

// read all records from the storer
func (s *storer[K, V]) ReadAll() ([]V, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data, nil
}

// read a record from the storer
func (s *storer[K, V]) ReadOne(id K) (*V, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.dataMap[id]
	if ok {

		return s.dataMap[id], nil
	}

	return nil, ErrorDataNotExists
}

type Writer[K comparable, V writer[K]] interface {
	Reader[K, V]
	writeFile() error

	Create(V) error
	Update(K, V) error
	Delete(K) error
}

type writer[K comparable] interface {
	GetID() K
	SetID(K)
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}
