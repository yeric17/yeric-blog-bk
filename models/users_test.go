package models

import (
	"testing"
)

func TestCreateAndDeleteUser(t *testing.T) {
	user := &User{}
	user.Name = "Test"
	user.Email = "correo@inventado.com"
	user.Password = "123456"

	err := user.Create()
	if err != nil {
		t.Errorf("Error creating user: %s", err)
	}

	if user.ID == "" {
		t.Errorf("User ID is empty")
	}

	if user.Name != "Test" {
		t.Errorf("User name is not Test")
	}

	if user.Email != "correo@inventado.com" {
		t.Errorf("User email is not same")
	}

	if user.Password != "123456" {
		t.Errorf("User password is not sames")
	}

	err = user.Delete()

	if err != nil {
		t.Errorf("Error deleting user: %s", err)
	}
}
