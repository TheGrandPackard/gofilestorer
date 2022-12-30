package filestorer

import (
	"sync"
	"time"

	"github.com/spf13/afero"
)

type Storer[D data] interface {
	Create(D) error
	Read() ([]D, error)
	Update(D) error
	Delete(uint64) error

	readFile() error
	writeFile() error
}

type storer[D data] struct {
	fs       afero.Fs
	fileName string
	mutex    sync.RWMutex
	data     []D
}

type data interface {
	GetID() uint64
	SetID(uint64)

	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
}
