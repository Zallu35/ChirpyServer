package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	id := uuid.New()
	_, err := MakeJWT(id, "ILikeCereal")
	if err != nil {
		t.Errorf("Error making JWT: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	id := uuid.New()
	ss, err := MakeJWT(id, "WhatSongLyricShouldIUseToday")
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

	/* test no longer needed as JWT duration is now fixed
	timeOut, err := MakeJWT(id, "ThisIsSupposedToFail")
	if err != nil {
		t.Errorf("Error making JWT for ValidateJWT timeout: %v", err)
	}

	toc, er := ValidateJWT(timeOut, "ThisIsSupposedToFail")
	if toc != uuid.Nil || er == nil {
		t.Errorf("JWT did not properly time out, validate failure")
	}*/
}
