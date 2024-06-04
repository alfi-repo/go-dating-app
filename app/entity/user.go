package entity

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserFailedToSave      = errors.New("failed to save user")
	ErrUserFailedToGetID     = errors.New("failed to get user ID")
	ErrUserFailedToFind      = errors.New("failed to find user")
	ErrUserInvalidEmail      = errors.New("invalid email address")
	ErrUserInvalidPassword   = errors.New("invalid password")
	ErrUserPasswordHash      = errors.New("failed to hash password")
	ErrUserPasswordIncorrect = errors.New("incorrect password")
)

type User struct {
	ID          int
	Email       string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SuspendedAt sql.NullTime
}

// CheckPassword check password. Return true if password is correct.
func (u *User) CheckPassword(plainPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword)); err != nil {
		return false, ErrUserPasswordIncorrect
	}
	return true, nil
}

// OnSave update timestamp on save.
func (u *User) OnSave() error {
	timeNow := time.Now().UTC()

	// New data or update
	if u.ID == 0 {
		// Hash password with bcrypt.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrUserPasswordHash, err)
		}

		u.Password = string(hashedPassword)
		u.CreatedAt = timeNow
		u.UpdatedAt = timeNow
	} else {
		u.UpdatedAt = timeNow
	}

	return nil
}

// NewUser create new user.
func NewUser(email, password string) (User, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Validate data.
	if err := validate.Var(email, "required,email,max=100"); err != nil {
		return User{}, ErrUserInvalidEmail
	}
	if err := validate.Var(password, "required,min=6"); err != nil {
		return User{}, ErrUserInvalidPassword
	}

	return User{
		Email:       email,
		Password:    password,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		SuspendedAt: sql.NullTime{Valid: false},
	}, nil
}
