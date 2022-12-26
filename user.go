package snsgo

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	AddUser(context.Context, User) (*User, error)
	GetUsers(context.Context) ([]User, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	// Update(User) error
}

type UserService interface {
	RegisterUser(context.Context, CreateRequest) error
	GetUsers(context.Context) ([]User, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Gender int8

const (
	Male Gender = iota + 1
	Female
	GenderUndefined
)

func (g Gender) String() string {
	switch g {
	case Male:
		return "Male"
	case Female:
		return "Female"
	}

	return ""
}

type User struct {
	ID          uuid.UUID
	Username    string
	Password    string
	Name        string
	Age         int
	Email       string
	Gender      Gender
	DateOfBirth time.Time
}

func (u *User) SetID() {
	id := uuid.New()
	u.ID = id
}

func (u *User) GetID() uuid.UUID {
	return u.ID
}

func (u *User) GetNewID() uuid.UUID {
	u.SetID()
	return u.ID
}

func (u *User) SetUsername(username string) {
	u.Username = username
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) SetGender(gender Gender) {
	u.Gender = gender
}

func (u *User) GetGender() Gender {
	return u.Gender
}

func (u *User) SetDateOfBirth(dob time.Time) {
	u.DateOfBirth = dob
}

func (u *User) GetDateOfBirth() time.Time {
	return u.DateOfBirth
}
