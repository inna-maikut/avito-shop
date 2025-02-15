package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) (*TransactionRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &TransactionRepository{
		db: db,
	}, nil
}

func (r *TransactionRepository) GetByEmployee(ctx context.Context, employeeID int64) ([]model.Transaction, error) {
	var transactions []EmployeeTransaction

	q := `SELECT true as is_sender, t.receiver_id as counterparty_employee_id, e.username as counterparty_username, t.amount
		FROM transaction t
		INNER JOIN employee e on e.id = t.receiver_id
		WHERE t.sender_id = $1
		UNION
		SELECT false as is_sender, t.sender_id as counterparty_employee_id, e.username as counterparty_username, t.amount
		FROM transaction t
		INNER JOIN employee e on e.id = t.sender_id
		WHERE t.receiver_id = $1
	`

	err := r.db.SelectContext(ctx, &transactions, q, employeeID)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext: %w", err)
	}

	res := make([]model.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		res = append(res, model.Transaction{
			IsSender:               transaction.IsSender,
			CounterpartyEmployeeID: transaction.CounterpartyEmployeeID,
			CounterpartyUsername:   transaction.CounterpartyUsername,
			Amount:                 transaction.Amount,
		})
	}

	return res, nil
}
