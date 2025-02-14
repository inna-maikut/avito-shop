package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(db *sqlx.DB) (*EmployeeRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &EmployeeRepository{
		db: db,
	}, nil
}

func (r *EmployeeRepository) GetByUsername(ctx context.Context, username string) (*model.Employee, error) {
	var employee Employee

	q := "SELECT id, username, password, balance FROM employee WHERE username = $1"

	err := r.db.GetContext(ctx, &employee, q, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("db.GetContext: %w", err)
	}

	return &model.Employee{
		ID:       employee.ID,
		Username: employee.Username,
		Password: employee.Password,
	}, nil
}

func (r *EmployeeRepository) Create(ctx context.Context, username, passwordHash string, balance int) (*model.Employee, error) {
	q := "INSERT INTO employee (username, password, balance) values " +
		"($1, $2, $3) " + // use binding to avoid SQL injection
		"ON CONFLICT DO NOTHING " +
		"RETURNING ID"

	var ID int64
	err := r.db.GetContext(ctx, &ID, q, username, passwordHash, balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrEmployeeAlreadyExists
		}
		return nil, fmt.Errorf("db.GetContext: %w", err)
	}

	return &model.Employee{
		ID:       ID,
		Username: username,
		Password: passwordHash,
		Balance:  balance,
	}, nil
}
