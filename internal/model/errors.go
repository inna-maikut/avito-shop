package model

import "errors"

var (
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrWrongEmployeePassword = errors.New("wrong employee password")
	ErrEmployeeAlreadyExists = errors.New("employee already exists")

	ErrMerchNotFound = errors.New("merch not found")

	ErrNotEnoughBalance               = errors.New("not enough balance")
	ErrSendingCoinsToMyselfNotAllowed = errors.New("sending coins to myself not allowed")
)
