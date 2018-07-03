package pica

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiVersionController_GetCommits(t *testing.T) {
	controller := NewApiVersionController("LICENSE")
	commits, err := controller.GetCommits()
	if err != nil {
		t.Error(err)
	}
	for _, item := range commits {
		t.Log(item)
	}
}

func TestApiVersionController_Notes(t *testing.T) {
	controller := NewApiVersionController("sample/pica.fun")
	notes, err := controller.Notes()
	if err != nil {
		t.Error(err)
	}
	t.Log(notes)
	assert.Equal(t, 1, len(notes.Changes))
}

func TestApiVersionController_Commit(t *testing.T) {
	controller := NewApiVersionController("LICENSE")
	hash, err := controller.Commit("test")
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
}
