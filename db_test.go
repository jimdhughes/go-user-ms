package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	DB = &DBClient{}
	TS = &TokenService{}
	DB.Initialize("./data/bolt-test.db")
	retCode := m.Run()
	os.Remove("./data/bolt-test.db")
	os.Exit(retCode)
}

func TestCheckUserShouldNotReturnError(t *testing.T) {
	email := "test@domain.com"
	isNew, err := DB.CheckUserIsNew(email)
	if err != nil {
		t.Errorf("An Error was thrown that should not have been: %s", err)
	}
	if isNew == false {
		t.Error("User shouldn't Exist but does")
	}
}

func TestCreateUser(t *testing.T) {
	u := User{
		Email:    "test@domain.com",
		Password: "password",
	}
	result, err := DB.CreateUser(u)
	if err != nil {
		t.Error(err)
	}
	if result == false {
		t.Error("User not created")
	}
}

func TestCreateUserWithSameEmailShouldFail(t *testing.T) {
	u := User{
		Email:    "test@domain.com",
		Password: "password",
	}
	result, err := DB.CreateUser(u)
	if err == nil {
		t.Error(err)
	}
	if result == true {
		t.Error("User not created")
	}
}

func TestCheckUserShouldReturnError(t *testing.T) {
	email := "test@domain.com"
	isNew, err := DB.CheckUserIsNew(email)
	if isNew == true {
		t.Error("User Should Exist but doesn't")
	}
	if err == nil {
		t.Error("Error should have been thrown")
	}
}

func TestValidLogin(t *testing.T) {
	_, err := DB.Login("test@domain.com", "password")
	if err != nil {
		t.Error(err)
	}
}

func TestInvalidLogin(t *testing.T) {
	_, err := DB.Login("test@domain.com", "wrongpassword")
	if err == nil {
		t.Error("The password should have been invalid")
	}
}

func TestGetUserByEmailInvalidEmail(t *testing.T) {
	user, err := DB.GetUserByEmail("invalidemail@somedomain.com")
	if err != nil {
		t.Errorf(err.Error())
	}
	if user.ID != "" {
		t.Errorf("No user should have been returned")
	}
}

func TestGetUserByEmailValidEmail(t *testing.T) {
	user, err := DB.GetUserByEmail("test@domain.com")
	if err != nil {
		t.Error(err.Error())
	}
	if user.Email == "" {
		t.Error("Expected to get user but got an empty user")
	}
	if user.Email != "test@domain.com" {
		t.Errorf("Returned a user other than expected")
	}
}
