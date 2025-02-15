package model

import "errors"

type Employee struct {
	ID       int64
	Username string
	Password string
	Balance  int64
}

var (
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrWrongEmployeePassword = errors.New("wrong employee password")
	ErrEmployeeAlreadyExists = errors.New("employee already exists")
)
