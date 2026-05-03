package main

import (
	"testing"
)

func TestCleanChirp(t *testing.T) {
	result := cleanChirp("I would like to visit Germany!")
	resultTwo := cleanChirp("I don't know what fornax means...")
	if result != "I would like to visit Germany!" {
		t.Errorf("Expected 'I would like to visit Germany!', Got '%v'", result)
	}
	if resultTwo != "I don't know what **** means..." {
		t.Errorf("Expected 'I don't know what **** means...', Got '%v'", resultTwo)
	}
}
