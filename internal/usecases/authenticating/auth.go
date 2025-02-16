package authenticating

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/inna-maikut/avito-shop/internal/model"
)

const initialCoins = 1000

type UseCase struct {
	employeeRepo  employeeRepo
	tokenProvider tokenProvider
}

func New(userRepo employeeRepo, tokenProvider tokenProvider) (*UseCase, error) {
	if userRepo == nil {
		return nil, errors.New("employeeRepo is nil")
	}
	if tokenProvider == nil {
		return nil, errors.New("tokenProvider is nil")
	}
	return &UseCase{
		employeeRepo:  userRepo,
		tokenProvider: tokenProvider,
	}, nil
}

func (uc *UseCase) Auth(ctx context.Context, username, password string) (string, error) {
	employee, err := uc.getOrCreateEmployee(ctx, username, password)
	if err != nil {
		if !errors.Is(err, model.ErrEmployeeAlreadyExists) {
			return "", fmt.Errorf("getOrCreateEmployee: %w", err)
		}

		// if user not found by username but has conflict on insert
		// we have concurrent auth request
		// repeat one time to get user by username and check password
		employee, err = uc.getOrCreateEmployee(ctx, username, password)
		if err != nil {
			return "", fmt.Errorf("getOrCreateEmployee retry: %w", err)
		}
	}

	token, err := uc.tokenProvider.CreateToken(employee.Username, employee.ID)
	if err != nil {
		return "", fmt.Errorf("tokenProvider.CreateToken: %w", err)
	}

	return token, nil
}

func (uc *UseCase) getOrCreateEmployee(ctx context.Context, username, password string) (*model.Employee, error) {
	employee, err := uc.employeeRepo.GetByUsername(ctx, username)
	if err == nil {
		err = uc.checkEmployeePassword(employee.Password, password)
		if err != nil {
			return nil, fmt.Errorf("checkEmployeePassword: %w", err)
		}

		return employee, nil
	}

	if !errors.Is(err, model.ErrEmployeeNotFound) {
		return nil, fmt.Errorf("employeeRepo.GetByUsername: %w", err)
	}

	employee, err = uc.createEmployee(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("createEmployee: %w", err)
	}

	return employee, nil
}

func (uc *UseCase) createEmployee(ctx context.Context, username, password string) (*model.Employee, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	employee, err := uc.employeeRepo.Create(ctx, username, string(hashedPassword), initialCoins)
	if err != nil {
		return nil, fmt.Errorf("employeeRepo.Create: %w", err)
	}

	return employee, nil
}

func (uc *UseCase) checkEmployeePassword(dbPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return model.ErrWrongEmployeePassword
		}

		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
	}

	return nil
}
