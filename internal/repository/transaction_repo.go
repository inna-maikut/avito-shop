package repository

import (
	"context"
	"errors"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type TransactionRepository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewTransactionRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) (*TransactionRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if getter == nil {
		return nil, errors.New("getter is nil")
	}

	return &TransactionRepository{
		db:     db,
		getter: getter,
	}, nil
}

func (r *TransactionRepository) trOrDB(ctx context.Context) trmsqlx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
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

	err := r.trOrDB(ctx).SelectContext(ctx, &transactions, q, employeeID)
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

func (r *TransactionRepository) Add(ctx context.Context, senderID, receiverID, amount int64) error {
	q := "INSERT INTO transaction (sender_id, receiver_id, amount) VALUES ($1, $2, $3)"

	_, err := r.trOrDB(ctx).ExecContext(ctx, q, senderID, receiverID, amount)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	return nil
}
