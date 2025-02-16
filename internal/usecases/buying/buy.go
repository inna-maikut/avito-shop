package buying

import (
	"context"
	"errors"
	"fmt"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type UseCase struct {
	trManager     trManager
	employeeRepo  employeeRepo
	inventoryRepo inventoryRepo
	merchRepo     merchRepo
}

func New(
	trManager trManager,
	employeeRepo employeeRepo,
	inventoryRepo inventoryRepo,
	merchRepo merchRepo,
) (*UseCase, error) {
	if trManager == nil {
		return nil, errors.New("trManager is nil")
	}
	if employeeRepo == nil {
		return nil, errors.New("employeeRepo is nil")
	}
	if inventoryRepo == nil {
		return nil, errors.New("inventoryRepo is nil")
	}
	if merchRepo == nil {
		return nil, errors.New("merchRepo is nil")
	}

	return &UseCase{
		trManager:     trManager,
		employeeRepo:  employeeRepo,
		inventoryRepo: inventoryRepo,
		merchRepo:     merchRepo,
	}, nil
}

func (uc *UseCase) Buy(ctx context.Context, employeeID int64, merchName string) error {
	merch, err := uc.merchRepo.GetByName(ctx, merchName)
	if err != nil {
		return fmt.Errorf("merchRepo.GetByName: %w", err)
	}

	err = uc.trManager.Do(ctx, func(ctx context.Context) (err error) {
		employee, err := uc.employeeRepo.GetByIDWithLock(ctx, employeeID)
		if err != nil {
			return fmt.Errorf("employeeRepo.GetByIDWithLock: %w", err)
		}

		if employee.Balance < merch.Price {
			return model.ErrNotEnoughBalance
		}

		err = uc.employeeRepo.IncreaseBalance(ctx, employeeID, -merch.Price)
		if err != nil {
			return fmt.Errorf("increase balance of current user with negative amount: %w", err)
		}

		err = uc.inventoryRepo.AddOne(ctx, employeeID, merch.ID)
		if err != nil {
			return fmt.Errorf("inventoryRepo.AddOne: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("trManager.Do: %w", err)
	}

	return nil
}
