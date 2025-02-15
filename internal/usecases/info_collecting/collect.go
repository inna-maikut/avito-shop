package info_collecting

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type UseCase struct {
	employeeRepo    employeeRepo
	transactionRepo transactionRepo
	inventoryRepo   inventoryRepo
}

func New(
	employeeRepo employeeRepo,
	transactionRepo transactionRepo,
	inventoryRepo inventoryRepo,
) (*UseCase, error) {
	if employeeRepo == nil {
		return nil, errors.New("employeeRepo is nil")
	}
	if transactionRepo == nil {
		return nil, errors.New("transactionRepo is nil")
	}
	if inventoryRepo == nil {
		return nil, errors.New("inventoryRepo is nil")
	}
	return &UseCase{
		employeeRepo:    employeeRepo,
		transactionRepo: transactionRepo,
		inventoryRepo:   inventoryRepo,
	}, nil
}

func (uc *UseCase) Collect(ctx context.Context, employeeID int64) (model.EmployeeInfo, error) {
	var (
		eg           *errgroup.Group
		employee     *model.Employee
		transactions []model.Transaction
		inventories  []model.Inventory
	)
	eg, ctx = errgroup.WithContext(ctx)

	eg.Go(func() (err error) {
		employee, err = uc.employeeRepo.GetByID(ctx, employeeID)
		if err != nil {
			return fmt.Errorf("employeeRepo.GetByID: %w", err)
		}

		return nil
	})

	eg.Go(func() (err error) {
		transactions, err = uc.transactionRepo.GetByEmployee(ctx, employeeID)
		if err != nil {
			return fmt.Errorf("transactionRepo.GetByEmployee: %w", err)
		}

		return nil
	})

	eg.Go(func() (err error) {
		inventories, err = uc.inventoryRepo.GetByEmployee(ctx, employeeID)
		if err != nil {
			return fmt.Errorf("inventoryRepo.GetByEmployee: %w", err)
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		return model.EmployeeInfo{}, fmt.Errorf("errgroup.Wait: %w", err)
	}

	var info model.EmployeeInfo
	if employee != nil { // err was nil, so employee is always not nil
		info.Coins = employee.Balance
	}

	info.Inventory = inventories

	info.SentTransactions = make([]model.Transaction, 0, len(transactions))
	info.ReceivedTransactions = make([]model.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		if transaction.IsSender {
			info.SentTransactions = append(info.SentTransactions, transaction)
		} else {
			info.ReceivedTransactions = append(info.ReceivedTransactions, transaction)
		}
	}

	return info, nil
}
