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
	Upsert(D) error

	readFile() error
	writeFile() error
}

type storer[D data] struct {
	fs       afero.Fs
	fileName string
	mutex    sync.RWMutex
	data     []D
}

type Data struct {
	ID        uint64    `json:"id" csv:"id"`
	CreatedAt time.Time `json:"created_at" csv:"created_at"`
}

type data interface {
	GetID() uint64
	SetID(uint64)

	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
}
