package coin_sending

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/inna-maikut/avito-shop/internal/model"
)

func TestUseCase_Buy(t *testing.T) {
	type mocks struct {
		trManager       *MocktrManager
		employeeRepo    *MockemployeeRepo
		transactionRepo *MocktransactionRepo
	}
	type args struct {
		employeeID     int64
		targetUsername string
		amount         int64
	}

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantErr error
	}{
		{
			name: "success.send",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(500)).
					Return(nil)
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(200)).
					Return(&model.Employee{
						ID:      200,
						Balance: 1000,
					}, nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(200), int64(-500)).
					Return(nil)

				m.transactionRepo.EXPECT().
					Add(gomock.Any(), int64(200), int64(100), int64(500)).
					Return(nil)
			},
			args: args{
				employeeID:     200,
				targetUsername: "test1",
				amount:         500,
			},
			wantErr: nil,
		},
		{
			name: "success.send.isTargetEmployeeIDGreaterThenSource",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(50)).
					Return(&model.Employee{
						ID:      50,
						Balance: 1000,
					}, nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(50), int64(-500)).
					Return(nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(500)).
					Return(nil)
				m.transactionRepo.EXPECT().
					Add(gomock.Any(), int64(50), int64(100), int64(500)).
					Return(nil)
			},
			args: args{
				employeeID:     50,
				targetUsername: "test1",
				amount:         500,
			},
			wantErr: nil,
		},
		{
			name: "error.SendingCoinsToMyselfNotAllowed",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       200,
						Username: "test1",
						Balance:  300,
					}, nil)
			},
			args: args{
				employeeID:     200,
				targetUsername: "test1",
				amount:         500,
			},
			wantErr: model.ErrSendingCoinsToMyselfNotAllowed,
		},
		{
			name: "error.transactionRepo.Add",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(500)).
					Return(nil)
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(200)).
					Return(&model.Employee{
						ID:      200,
						Balance: 1000,
					}, nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(200), int64(-500)).
					Return(nil)

				m.transactionRepo.EXPECT().
					Add(gomock.Any(), int64(200), int64(100), int64(500)).
					Return(assert.AnError)
			},
			args: args{
				employeeID:     200,
				targetUsername: "test1",
				amount:         500,
			},
			wantErr: assert.AnError,
		},
		{
			name: "error.employeeRepo.IncreaseBalance.isTargetEmployeeIDGreaterThenSource",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(50)).
					Return(&model.Employee{
						ID:      50,
						Balance: 1000,
					}, nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(50), int64(-500)).
					Return(nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(500)).
					Return(assert.AnError)
			},
			args: args{
				employeeID:     50,
				targetUsername: "test1",
				amount:         500,
			},
			wantErr: assert.AnError,
		},
		{
			name: "error.employeeRepo.IncreaseBalanc",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(500)).
					Return(nil)
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(200)).
					Return(&model.Employee{
						ID:      200,
						Balance: 1000,
					}, nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(200), int64(-500)).
					Return(assert.AnError)
			},
			args: args{
				employeeID:     200,
				targetUsername: "test1",
				amount:         500,
			},
			wantErr: assert.AnError,
		},
		{
			name: "error.employeeRepo.GetByIDWithLock",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(500)).
					Return(nil)
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(200)).
					Return(nil, assert.AnError)
			},
			args: args{
				employeeID:     200,
				targetUsername: "test1",
				amount:         500,
			},
			wantErr: assert.AnError,
		},
		{
			name: "error.employeeRepo.IncreaseBalance.!isTargetEmployeeIDGreaterThenSource",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(500)).
					Return(assert.AnError)
			},
			args: args{
				employeeID:     200,
				targetUsername: "test1",
				amount:         500,
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				employeeRepo:    NewMockemployeeRepo(ctrl),
				trManager:       NewMocktrManager(ctrl),
				transactionRepo: NewMocktransactionRepo(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.trManager, m.employeeRepo, m.transactionRepo)
			require.NoError(t, err)

			err = uc.Send(context.Background(), tc.args.employeeID, tc.args.targetUsername, tc.args.amount)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
