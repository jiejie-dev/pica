package pica

import "testing"

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
	controller := NewApiVersionController("LICENSE")
	notes, err := controller.Notes()
	if err != nil {
		t.Error(err)
	}
	t.Log(notes)
}

func TestApiVersionController_Commit(t *testing.T) {
	controller := NewApiVersionController("LICENSE")
	hash, err := controller.Commit("test")
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
}
