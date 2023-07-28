package domain

import "errors"

var (
	ErrUserNotFound                  = errors.New("user not found")
	ErrUserWithSameLoginAlreadyExist = errors.New("user with the same login already exist")
)

type User struct {
	Id        uint
	Login     string
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

type UserParameters struct {
	Login     string
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

