package filestorer

import (
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type testJSONData struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

func (d *testJSONData) GetID() uint64 {
	return d.ID
}

func (d *testJSONData) SetID(id uint64) {
	d.ID = id
}

func (d *testJSONData) SetCreatedAt(createdAt time.Time) {
	d.CreatedAt = createdAt
}

func getJSONFilesystem(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	// create test files and directories
	err := fs.MkdirAll("data", 0755)
	assert.NoError(t, err)
	err = afero.WriteFile(fs, "data.json", []byte(`[
		{
			"id": 1,
			"created_at": "2022-12-27T12:45:51.8347046-08:00",
			"name": "Foobar"
		}
	]`), 0644)
	assert.NoError(t, err)
	err = afero.WriteFile(fs, "invalid.json", []byte(``), 0644)
	assert.NoError(t, err)

	return fs
}

func TestJSONReader(t *testing.T) {
	fs := getJSONFilesystem(t)

	// Read non-existant file
	s, err := NewJSONReader[*testJSONData](fs, "./foobar.json")
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read invalid file
	s, err = NewJSONReader[*testJSONData](fs, "./data/invalid.json")
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read test file
	s, err = NewJSONReader[*testJSONData](fs, "./data.json")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Read
	read, err := s.Read()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 1)
	assert.Equal(t, uint64(1), read[0].ID)
	assert.Equal(t, "Foobar", read[0].Name)
	assert.NotEmpty(t, read[0].CreatedAt)
}

func TestJSONWriter(t *testing.T) {
	fs := getJSONFilesystem(t)

	// Read non-existant file
	s, err := NewJSONWriter[*testJSONData](fs, "./foobar.json")
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read invalid file
	s, err = NewJSONWriter[*testJSONData](fs, "./data/invalid.json")
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read test file
	s, err = NewJSONWriter[*testJSONData](fs, "./data.json")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Read
	read, err := s.Read()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 1)
	assert.Equal(t, uint64(1), read[0].ID)
	assert.Equal(t, "Foobar", read[0].Name)
	assert.NotEmpty(t, read[0].CreatedAt)

	// Create
	data := &testJSONData{Name: "new"}
	err = s.Create(data)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), data.ID)
	assert.Equal(t, "new", data.Name)
	assert.NotEmpty(t, data.CreatedAt)

	read, err = s.Read()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 2)
	assert.Equal(t, uint64(2), read[1].ID)
	assert.Equal(t, "new", read[1].Name)
	assert.NotEmpty(t, read[1].CreatedAt)

	// Upsert - Update
	data.Name = "upserted"
	err = s.Upsert(data)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), data.ID)
	assert.Equal(t, "upserted", data.Name)
	assert.NotEmpty(t, data.CreatedAt)

	read, err = s.Read()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 2)
	assert.Equal(t, uint64(2), read[1].ID)
	assert.Equal(t, "upserted", read[1].Name)
	assert.NotEmpty(t, read[1].CreatedAt)

	// Upsert - Insert
	upsert := &testJSONData{Name: "upsert"}
	err = s.Upsert(upsert)
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), upsert.ID)
	assert.Equal(t, "upsert", upsert.Name)
	assert.NotEmpty(t, upsert.CreatedAt)

	read, err = s.Read()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 3)
	assert.Equal(t, uint64(3), read[2].ID)
	assert.Equal(t, "upsert", read[2].Name)
	assert.NotEmpty(t, read[2].CreatedAt)

	// Update
	data.Name = "updated"
	err = s.Update(data)
	assert.NoError(t, err)

	read, err = s.Read()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 3)
	assert.Equal(t, uint64(2), read[1].ID)
	assert.Equal(t, "updated", read[1].Name)
	assert.NotEmpty(t, read[1].CreatedAt)

	// Delete
	err = s.Delete(data.ID)
	assert.NoError(t, err)

	read, err = s.Read()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 2)

	// Update - Not Exists
	err = s.Update(data)
	assert.Error(t, err)

	// Delete - Not Exists
	err = s.Delete(data.ID)
	assert.Error(t, err)
}
