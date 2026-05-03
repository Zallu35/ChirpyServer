package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	id := uuid.New()
	_, err := MakeJWT(id, "ILikeCereal", 30*time.Second)
	if err != nil {
		t.Errorf("Error making JWT: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	id := uuid.New()
	ss, err := MakeJWT(id, "WhatSongLyricShouldIUseToday", 5*time.Minute)
	if err != nil {
		t.Errorf("Error making JWT for ValidateJWT test: %v", err)
	}
	check, er := ValidateJWT(ss, "WhatSongLyricShouldIUseToday")
	if er != nil {
		t.Errorf("ValidateJWT failure: %v", er)
	}
	if check != id {
		t.Errorf("ValidateJWT: IDs do not match!\nInitial: %v\nOutput: %v\n", id, check)
	}

	timeOut, err := MakeJWT(id, "ThisIsSupposedToFail", -2*time.Second)
	if err != nil {
		t.Errorf("Error making JWT for ValidateJWT timeout: %v", err)
	}

	toc, er := ValidateJWT(timeOut, "ThisIsSupposedToFail")
	if toc != uuid.Nil || er == nil {
		t.Errorf("JWT did not properly time out, validate failure")
	}
}
