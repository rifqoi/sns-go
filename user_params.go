package snsgo

import (
	"errors"
	"time"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var (
	ErrDateNotValid   = errors.New("Date not valid, date format must be YYYY-MM-dd")
	ErrGenderNotValid = errors.New("Gender must be either male or female.")
	ErrPassword       = errors.New("Password must be minimum eight characters, at least one uppercase letter, one lowercase letter, one number and one special character")
)

type CreateRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Gender      Gender `json:"gender"`
	Age         int    `json:"age"`
	DateOfBirth string `json:"date_of_birth"`
}

func (r CreateRequest) NewUser() User {
	t, _ := time.Parse("2006-01-02", r.DateOfBirth)
	user := User{
		Username:    r.Username,
		Name:        r.Name,
		Email:       r.Email,
		Password:    r.Password,
		Age:         r.Age,
		Gender:      r.Gender,
		DateOfBirth: t,
	}

	return user
}

func (r CreateRequest) Validate() error {

	username := validation.Field(&r.Username, validation.Required, validation.Length(7, 15), is.Alphanumeric)
	name := validation.Field(&r.Name, validation.Required)
	password := validation.Field(
		&r.Password,
		validation.Required,
		validation.Length(8, 30),
		validation.By(checkPassword),
	)
	age := validation.Field(&r.Age, validation.Required, validation.Min(8), validation.Max(30))
	gender := validation.Field(&r.Gender, validation.Required, validation.In(Male, Female).Error("either be Female or Male"))
	dob := validation.Field(&r.DateOfBirth, validation.Required, validation.Date("2006-01-02"))
	email := validation.Field(&r.Email, validation.Required, is.Email)

	err := validation.ValidateStruct(&r,
		username,
		email,
		name,
		age,
		password,
		gender,
		dob,
	)

	return err
}

type UpdateUser struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Gender      Gender `json:"gender"`
	Age         int    `json:"age"`
	DateOfBirth string `json:"date_of_birth"`
}

func (u *UpdateUser) NewUser() User {
	t, _ := time.Parse("2006-01-02", u.DateOfBirth)
	return User{
		Name:        u.Name,
		Password:    u.Password,
		Gender:      u.Gender,
		Age:         u.Age,
		DateOfBirth: t,
	}
}

func (u *UpdateUser) Validate() error {
	name := validation.Field(&u.Name, validation.Required)
	password := validation.Field(
		&u.Password,
		validation.Required,
		validation.Length(8, 30),
		validation.By(checkPassword),
	)
	age := validation.Field(&u.Age, validation.Required, validation.Min(8), validation.Max(30))
	gender := validation.Field(&u.Gender, validation.Required, validation.In(Male, Female).Error("either be Female or Male"))
	dob := validation.Field(&u.DateOfBirth, validation.Required, validation.Date("2006-01-02"))

	err := validation.ValidateStruct(&u,
		name,
		age,
		password,
		gender,
		dob,
	)

	return err
}

func checkPassword(value interface{}) error {
	s, _ := value.(string)
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if hasUpper && hasLower && hasNumber && hasSpecial {
		return nil
	}

	return ErrPassword
}
