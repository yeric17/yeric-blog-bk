package models

import "testing"

func TestCreateAndDeleteUser(t *testing.T) {
	user := User{
		Name:     "Yerik",
		Email:    "correo@inventado.com",
		Password: "123456",
		Status:   "active",
		RoleID:   1,
	}
	err := user.Create()

	if err != nil {
		t.Errorf("Error creating user: %s", err)
	}

	if user.ID == "" {
		t.Errorf("User ID is empty")
	}

	if user.Name != "Yerik" {

		t.Errorf("User name is not correct")
	}

	if user.Email != "correo@inventado.com" {
		t.Errorf("User email is not correct")
	}

	if user.Status != "active" {
		t.Errorf("User status is not correct")
	}

	if user.RoleID != 1 {
		t.Errorf("User role ID is not correct")
	}

	err = user.Delete()

	if err != nil {
		t.Errorf("Error deleting user: %s", err)
	}
}
