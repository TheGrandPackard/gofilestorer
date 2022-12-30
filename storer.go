package filestorer

import (
	"sync"
	"time"

	"github.com/spf13/afero"
)

type Reader[D any] interface {
	readFile() error

	Read() ([]D, error)
}

type Writer[D any] interface {
	readFile() error
	writeFile() error

	Create(D) error
	Read() ([]D, error)
	Update(D) error
	Delete(uint64) error
	Upsert(D) error
}

type storer[D any] struct {
	fs       afero.Fs
	fileName string
	mutex    sync.RWMutex
	data     []D
}

type Data struct {
	ID        uint64    `json:"id" csv:"id"`
	CreatedAt time.Time `json:"created_at" csv:"created_at"`
}

type reader interface {
	GetID() uint64
}

type writer interface {
	GetID() uint64
	SetID(uint64)
	SetCreatedAt(time.Time)
}
