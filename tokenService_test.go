package main

import (
	"testing"
)

func TestCreateTokenFromNotValidUser(t *testing.T) {
	user := User{}
	token, err := TS.CreateToken(user)
	if err == nil {
		t.Errorf("Token should not be created for an empty user")
	}
	if token != "" {
		t.Errorf("Error should have been thrown as token cannot be created")
	}
}

func TestCreateTokenFromValidUser(t *testing.T) {
	user := User{Email: "tokentester@domain.com", Password: "password"}
	DB.CreateUser(user)                     // we test this elsewhere so assume it works. see db_test.go
	u, err := DB.GetUserByEmail(user.Email) //we can ignore the error as we test this elsewhere
	if err != nil {
		t.Errorf("%s", err)
	}
	token, err := TS.CreateToken(u)
	if err != nil {
		t.Errorf("Token could not be created for user: %s", err)
	}
	if token == "" {
		t.Errorf("Empty token was returned. Expected a token as a string")
	}
}

func TestValidateTokenFromInvalidToken(t *testing.T) {
	token := "thisisNOTatoken"
	payload, err := TS.ValidateToken(token)
	if err == nil {
		t.Errorf("Expected an error to be thrown as token is not valid")
	}
	if payload != nil {
		t.Errorf("Expected no user to be returned")
	}
}

func TestValidateTokenFromValidToken(t *testing.T) {
	user, err := DB.GetUserByEmail("tokentester@domain.com")
	if err != nil {
		t.Error(err)
	}
	if user.ID == "" {
		t.Errorf("Expected to get a user. Got none")
	}
	token, err := TS.CreateToken(user)
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("Got an empty token when expecting a real token")
	}
	payload, err := TS.ValidateToken(token)
	if err != nil {
		t.Error(err)
	}
	if payload == nil {
		t.Error("Expected to get a payload, got none")
	}
}
