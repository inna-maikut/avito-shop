package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type EmployeeRepository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewEmployeeRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) (*EmployeeRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if getter == nil {
		return nil, errors.New("getter is nil")
	}

	return &EmployeeRepository{
		db:     db,
		getter: getter,
	}, nil
}

func (r *EmployeeRepository) trOrDB(ctx context.Context) trmsqlx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
}

func (r *EmployeeRepository) GetByUsername(ctx context.Context, username string) (*model.Employee, error) {
	var employee Employee

	q := "SELECT id, username, password, balance FROM employee WHERE username = $1"

	err := r.trOrDB(ctx).GetContext(ctx, &employee, q, username)
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
		Balance:  employee.Balance,
	}, nil
}

func (r *EmployeeRepository) GetByID(ctx context.Context, employeeID int64) (*model.Employee, error) {
	var employee Employee

	q := "SELECT id, username, password, balance FROM employee WHERE id = $1"

	err := r.trOrDB(ctx).GetContext(ctx, &employee, q, employeeID)
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
		Balance:  employee.Balance,
	}, nil
}

func (r *EmployeeRepository) Create(ctx context.Context, username, passwordHash string, balance int64) (*model.Employee, error) {
	q := "INSERT INTO employee (username, password, balance) values " +
		"($1, $2, $3) " + // use binding to avoid SQL injection
		"ON CONFLICT DO NOTHING " +
		"RETURNING ID"

	var ID int64
	err := r.trOrDB(ctx).GetContext(ctx, &ID, q, username, passwordHash, balance)
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

func (r *EmployeeRepository) GetByIDWithLock(ctx context.Context, employeeID int64) (*model.Employee, error) {
	var employee Employee

	q := "SELECT id, username, password, balance FROM employee WHERE id = $1 FOR NO KEY UPDATE"

	err := r.trOrDB(ctx).GetContext(ctx, &employee, q, employeeID)
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
		Balance:  employee.Balance,
	}, nil
}

func (r *EmployeeRepository) IncreaseBalance(ctx context.Context, employeeID, amount int64) error {
	q := "UPDATE employee SET balance = balance + $2 WHERE id = $1"

	_, err := r.trOrDB(ctx).ExecContext(ctx, q, employeeID, amount)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}
