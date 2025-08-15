package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID     string
	Name       string
	PW         string
	DebutDate  time.Time // For calculating experience
	PlayerType float32   // 1-5:General-Anison
	Stance     float32   // 1-5:Gachi-Enjoy
	Club       string
	Major      string
	Profile    string
}

type UserUpdateRequest struct {
	UserID string
	// PW         string
	Name       *string
	DebutDate  *time.Time // For calculating experience
	PlayerType *float32   // 1-5:General-Anison
	Stance     *float32   // 1-5:Gachi-Enjoy
	Club       *string
	Major      *string
	Profile    *string
}

var (
	ErrAlreadyExistUser = errors.New("already exist user")
	ErrUserNotFound     = errors.New("user not found")
	ErrInternalDBError  = errors.New("internal database error")
	ErrInvalidError     = errors.New("database invalid error raised")
)

func (r Database) CreateUser(u *User) (*User, error) {
	_, err := r.GetUser(u.UserID)
	if err == nil {
		return nil, ErrAlreadyExistUser
	} else if err != ErrUserNotFound {
		return nil, ErrInvalidError
	}

	result := r.db.Create(u)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.GetUser(u.UserID)
}

func (r Database) GetUser(userID string) (*User, error) {
	var users []User
	tx := r.db.Where("user_id = ?", userID).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	if len(users) != 1 {
		return nil, ErrInternalDBError
	}

	return &users[0], nil
}

func (r Database) UpdateUser(u *User) (*User, error) {
	return nil, nil
}

func (r Database) DeleteUser(userID string) error {
	return nil
}
