package pkg

import "testing"

func TestExists(t *testing.T) {
	// Since this may varies depends on OS, i had tested
	// only with this modules's folder instead of custom folder.
	var dummyFolder = []string{
		"../cmd",
		"../configs",
		"../pkg",
		"../temp",
	}

	for _, path := range dummyFolder {
		if !Exists(path) {
			t.Errorf("Folder %s doest not exist", path)
		}
	}
}

func TestContains(t *testing.T) {
	var strings = []string{
		"hello",
		"world",
		"from",
		"me",
	}

	for _, s:= range strings {
		if !Contains(strings, s) {
			t.Errorf("No wording matched with %s", s)
		}
	}
}