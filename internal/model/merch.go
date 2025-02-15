package model

import "errors"

type Merch struct {
	ID    int64
	Name  string
	Price int64
}

var ErrMerchNotFound = errors.New("merch not found")
