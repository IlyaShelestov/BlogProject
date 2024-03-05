package mocks

import (
	"blog_project/internal/models"
	"time"
)

type UserModel struct{}

func (m *UserModel) Insert(username, password string) error {
	switch username {
	case "dupedupe":
		return models.ErrDuplicateUsername
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(username, password string) (int, error) {
	if username == "alicealice" && password == "pa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) ExistsByUsername(username string) (bool, error) {
	switch username {
	case "dupedupe":
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (models.User, error) {
	if id == 1 {
		u := models.User{
			ID:       1,
			Username: "Alice",
			Created:  time.Now(),
		}
		return u, nil
	}
	return models.User{}, models.ErrNoRecord
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword != "pa$$word" {
			return models.ErrInvalidCredentials
		}
		return nil
	}
	return models.ErrNoRecord
}

func (m *UserModel) IsAdmin(id int) (bool, error) {
	return true, nil
}
