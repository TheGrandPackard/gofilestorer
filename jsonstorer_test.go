package gofilestorer

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type testJSONDataUUID struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Name      string     `json:"name"`
}

func (d *testJSONDataUUID) GetID() uuid.UUID {
	return d.ID
}

func (d *testJSONDataUUID) SetID(id uuid.UUID) {
	d.ID = id
}

func (d *testJSONDataUUID) SetCreatedAt(createdAt time.Time) {
	d.CreatedAt = createdAt
}

func (d *testJSONDataUUID) SetUpdatedAt(updatedAt time.Time) {
	d.UpdatedAt = &updatedAt
}

func getJSONFilesystem(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	// create test files and directories
	err := fs.MkdirAll("data", 0755)
	assert.NoError(t, err)
	err = afero.WriteFile(fs, "uuid.json", []byte(`[
		{
			"id": "e21ab9b3-bb4e-4921-815b-41de7980c5da",
			"created_at": "2022-12-27T12:45:51.8347046-08:00",
			"name": "Foobar"
		}
	]`), 0644)
	assert.NoError(t, err)
	err = afero.WriteFile(fs, "int64.json", []byte(`[
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
	s, err := NewJSONReader[uuid.UUID, *testJSONDataUUID](fs, "./foobar.json")
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read invalid file
	s, err = NewJSONReader[uuid.UUID, *testJSONDataUUID](fs, "./data/invalid.json")
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read test file
	s, err = NewJSONReader[uuid.UUID, *testJSONDataUUID](fs, "./uuid.json")
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Read
	read, err := s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 1)
	assert.Equal(t, "Foobar", read[0].Name)
	assert.NotEmpty(t, read[0].ID)
	assert.NotEmpty(t, read[0].CreatedAt)
}

func TestJSONWriterUUID(t *testing.T) {
	fs := getJSONFilesystem(t)

	newIdFunc := func(data []*testJSONDataUUID) uuid.UUID {
		return uuid.New()
	}

	// Read non-existant file
	s, err := NewJSONWriter[uuid.UUID, *testJSONDataUUID](fs, "./foobar.json", newIdFunc)
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read invalid file
	s, err = NewJSONWriter[uuid.UUID, *testJSONDataUUID](fs, "./data/invalid.json", newIdFunc)
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read test file
	s, err = NewJSONWriter[uuid.UUID, *testJSONDataUUID](fs, "./uuid.json", newIdFunc)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Read
	read, err := s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 1)
	assert.Equal(t, "e21ab9b3-bb4e-4921-815b-41de7980c5da", read[0].ID.String())
	assert.Equal(t, "Foobar", read[0].Name)
	assert.NotEmpty(t, read[0].CreatedAt)

	// Create
	data := &testJSONDataUUID{Name: "new"}
	_, err = s.Create(data)
	assert.NoError(t, err)
	assert.NoError(t, uuid.Validate(data.ID.String()))
	assert.Equal(t, "new", data.Name)
	assert.NotEmpty(t, data.CreatedAt)

	read, err = s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 2)
	assert.NoError(t, uuid.Validate(data.ID.String()))
	assert.Equal(t, "new", read[1].Name)
	assert.NotEmpty(t, read[1].CreatedAt)

	// Update
	data.Name = "updated"
	_, err = s.Update(data.GetID(), data)
	assert.NoError(t, err)

	read, err = s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 2)
	assert.NoError(t, uuid.Validate(data.ID.String()))
	assert.Equal(t, "updated", read[1].Name)
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

type testJSONDataInt64 struct {
	ID        int64      `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Name      string     `json:"name"`
}

func (d *testJSONDataInt64) GetID() int64 {
	return d.ID
}

func (d *testJSONDataInt64) SetID(id int64) {
	d.ID = id
}

func (d *testJSONDataInt64) SetCreatedAt(createdAt time.Time) {
	d.CreatedAt = createdAt
}

func (d *testJSONDataInt64) SetUpdatedAt(updatedAt time.Time) {
	d.UpdatedAt = &updatedAt
}

func TestJSONWriterInt64(t *testing.T) {
	fs := getJSONFilesystem(t)

	newIdFunc := func(data []*testJSONDataInt64) int64 {
		return int64(len(data) + 1)
	}

	// Read non-existant file
	s, err := NewJSONWriter[int64, *testJSONDataInt64](fs, "./foobar.json", newIdFunc)
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read invalid file
	s, err = NewJSONWriter[int64, *testJSONDataInt64](fs, "./data/invalid.json", newIdFunc)
	assert.Error(t, err)
	assert.Nil(t, s)

	// Read test file
	s, err = NewJSONWriter[int64, *testJSONDataInt64](fs, "./int64.json", newIdFunc)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Read
	read, err := s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 1)
	assert.Equal(t, "Foobar", read[0].Name)
	assert.NotEmpty(t, int64(1), read[0].ID)
	assert.NotEmpty(t, read[0].CreatedAt)

	// Create
	data := &testJSONDataInt64{Name: "new"}
	_, err = s.Create(data)
	assert.NoError(t, err)
	assert.Equal(t, "new", data.Name)
	assert.Equal(t, int64(2), data.ID)
	assert.NotEmpty(t, data.CreatedAt)

	read, err = s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 2)
	assert.Equal(t, int64(2), read[1].ID)
	assert.Equal(t, "new", read[1].Name)
	assert.NotEmpty(t, read[1].CreatedAt)

	// Update
	data.Name = "updated"
	_, err = s.Update(data.GetID(), data)
	assert.NoError(t, err)

	read, err = s.ReadAll()
	assert.NoError(t, err)
	assert.NotNil(t, read)
	assert.Len(t, read, 2)
	assert.Equal(t, int64(2), read[1].ID)
	assert.Equal(t, "updated", read[1].Name)
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
