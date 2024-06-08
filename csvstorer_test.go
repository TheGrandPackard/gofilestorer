package gofilestorer

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type testCSVData struct {
	ID        uuid.UUID  `csv:"id"`
	CreatedAt time.Time  `csv:"created_at"`
	UpdatedAt *time.Time `csv:"updated_at"`
	Name      string     `csv:"name"`
}

func (d *testCSVData) GetID() uuid.UUID {
	return d.ID
}

func (d *testCSVData) SetID(id uuid.UUID) {
	d.ID = id
}

func (d *testCSVData) GetCreatedAt() time.Time {
	return d.CreatedAt
}

func (d *testCSVData) SetCreatedAt(createdAt time.Time) {
	d.CreatedAt = createdAt
}

func (d *testCSVData) SetUpdatedAt(updatedAt time.Time) {
	d.UpdatedAt = &updatedAt
}

func getCSVFilesystem(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	// create test files and directories
	err := fs.MkdirAll("data", 0755)
	assert.NoError(t, err)
	err = afero.WriteFile(fs, "data.json", []byte(`id;created_at;name
e21ab9b3-bb4e-4921-815b-41de7980c5da;2022-12-27T12:45:51.8347046-08:00;Foobar`), 0644)
	assert.NoError(t, err)
	err = afero.WriteFile(fs, "invalid.json", []byte(``), 0644)
	assert.NoError(t, err)

	return fs
}

func TestCSVReader(t *testing.T) {
	fs := getCSVFilesystem(t)

	// Read non-existant file
	s, err := NewCSVReader[uuid.UUID, *testCSVData](fs, "./foobar.json", ';')
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read invalid file
	s, err = NewCSVReader[uuid.UUID, *testCSVData](fs, "./data/invalid.json", ';')
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read test file
	s, err = NewCSVReader[uuid.UUID, *testCSVData](fs, "./data.json", ';')
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Read
	read, err := s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 1)
	assert.Equal(t, uuid.MustParse("e21ab9b3-bb4e-4921-815b-41de7980c5da"), read[0].ID)
	assert.Equal(t, "Foobar", read[0].Name)
	assert.NotEmpty(t, read[0].CreatedAt)
}

func TestCSVWriter(t *testing.T) {
	fs := getCSVFilesystem(t)

	newIdFunc := func(data []*testCSVData) uuid.UUID {
		return uuid.New()
	}

	// Read non-existant file
	s, err := NewCSVWriter[uuid.UUID, *testCSVData](fs, "./foobar.json", ';', newIdFunc)
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read invalid file
	s, err = NewCSVWriter[uuid.UUID, *testCSVData](fs, "./data/invalid.json", ';', newIdFunc)
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read test file
	s, err = NewCSVWriter[uuid.UUID, *testCSVData](fs, "./data.json", ';', newIdFunc)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Read
	read, err := s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 1)
	assert.Equal(t, uuid.MustParse("e21ab9b3-bb4e-4921-815b-41de7980c5da"), read[0].ID)
	assert.Equal(t, "Foobar", read[0].Name)
	assert.NotEmpty(t, read[0].CreatedAt)

	// Create
	data := &testCSVData{Name: "new"}
	_, err = s.Create(data)
	assert.NoError(t, err)
	assert.Equal(t, "new", data.Name)
	assert.NotEmpty(t, data.ID)
	assert.NotEmpty(t, data.CreatedAt)

	read, err = s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 2)
	assert.Equal(t, "new", read[1].Name)
	assert.NotEmpty(t, data.ID)
	assert.NotEmpty(t, read[1].CreatedAt)

	// Update
	data.Name = "updated"
	_, err = s.Update(data.GetID(), data)
	assert.NoError(t, err)

	read, err = s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 2)
	assert.Equal(t, "updated", read[1].Name)
	assert.NotEmpty(t, data.ID)
	assert.NotEmpty(t, read[1].CreatedAt)

	// Delete
	err = s.Delete(data.ID)
	assert.NoError(t, err)

	read, err = s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 1)

	// Update - Not Exists
	_, err = s.Update(data.GetID(), data)
	assert.Error(t, err)

	// Delete - Not Exists
	err = s.Delete(data.ID)
	assert.Error(t, err)
}
