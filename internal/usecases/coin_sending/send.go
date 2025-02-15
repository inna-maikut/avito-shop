package coin_sending

import (
	"context"
	"errors"
	"fmt"

	"github.com/inna-maikut/avito-shop/internal/model"
)

type UseCase struct {
	trManager       trManager
	employeeRepo    employeeRepo
	transactionRepo transactionRepo
}

func New(
	trManager trManager,
	employeeRepo employeeRepo,
	transactionRepo transactionRepo,
) (*UseCase, error) {
	if trManager == nil {
		return nil, errors.New("trManager is nil")
	}
	if employeeRepo == nil {
		return nil, errors.New("employeeRepo is nil")
	}
	if transactionRepo == nil {
		return nil, errors.New("transactionRepo is nil")
	}

	return &UseCase{
		trManager:       trManager,
		employeeRepo:    employeeRepo,
		transactionRepo: transactionRepo,
	}, nil
}

func (uc *UseCase) Send(ctx context.Context, employeeID int64, targetUsername string, amount int64) error {
	targetEmployee, err := uc.employeeRepo.GetByUsername(ctx, targetUsername)
	if err != nil {
		return fmt.Errorf("employeeRepo.GetByUsername: %w", err)
	}

	targetEmployeeID := targetEmployee.ID

	if targetEmployeeID == employeeID {
		return model.ErrSendingCoinsToMyselfNotAllowed
	}

	isTargetEmployeeIDGreaterThenSource := targetEmployeeID > employeeID // couldn't be equal because of the check above

	err = uc.trManager.Do(ctx, func(ctx context.Context) (err error) {
		// need to follow lock order to avoid deadlocks
		// first lock lower employeeID with either IncreaseBalance or GetByIDWithLock
		if !isTargetEmployeeIDGreaterThenSource {
			err = uc.employeeRepo.IncreaseBalance(ctx, targetEmployeeID, amount)
			if err != nil {
				return fmt.Errorf("increase balance of target employee with lower employeeID: %w", err)
			}
		}

		employee, err := uc.employeeRepo.GetByIDWithLock(ctx, employeeID)
		if err != nil {
			return fmt.Errorf("employeeRepo.GetByIDWithLock: %w", err)
		}

		if employee.Balance < amount {
			return model.ErrNotEnoughBalance
		}

		err = uc.employeeRepo.IncreaseBalance(ctx, employeeID, -amount)
		if err != nil {
			return fmt.Errorf("increase balance of current user with negative amount: %w", err)
		}

		if isTargetEmployeeIDGreaterThenSource {
			err = uc.employeeRepo.IncreaseBalance(ctx, targetEmployeeID, amount)
			if err != nil {
				return fmt.Errorf("increase balance of target employee with greater employeeID: %w", err)
			}
		}

		err = uc.transactionRepo.Add(ctx, employeeID, targetEmployeeID, amount)
		if err != nil {
			return fmt.Errorf("transactionRepo.Add: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("trManager.Do: %w", err)
	}

	return nil
}
