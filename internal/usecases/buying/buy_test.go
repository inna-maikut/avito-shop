package buying

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
		trManager     *MocktrManager
		employeeRepo  *MockemployeeRepo
		inventoryRepo *MockinventoryRepo
		merchRepo     *MockmerchRepo
	}
	type args struct {
		employeeID int64
		merchName  string
	}

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantErr error
	}{
		{
			name: "success.buy",
			prepare: func(m *mocks) {
				m.merchRepo.EXPECT().
					GetByName(gomock.Any(), "test1").
					Return(&model.Merch{
						ID:    1,
						Name:  "test1",
						Price: 300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(100)).
					Return(&model.Employee{
						ID:      100,
						Balance: 1000,
					}, nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(-300)).
					Return(nil)
				m.inventoryRepo.EXPECT().
					AddOne(gomock.Any(), int64(100), int64(1)).
					Return(nil)
			},
			args: args{
				employeeID: 100,
				merchName:  "test1",
			},
			wantErr: nil,
		},
		{
			name: "error.NotEnoughBalance",
			prepare: func(m *mocks) {
				m.merchRepo.EXPECT().
					GetByName(gomock.Any(), "test1").
					Return(&model.Merch{
						ID:    1,
						Name:  "test1",
						Price: 3000,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(100)).
					Return(&model.Employee{
						ID:      100,
						Balance: 1000,
					}, nil)
			},
			args: args{
				employeeID: 100,
				merchName:  "test1",
			},
			wantErr: model.ErrNotEnoughBalance,
		},
		{
			name: "error.merchRepo.GetByName",
			prepare: func(m *mocks) {
				m.merchRepo.EXPECT().
					GetByName(gomock.Any(), "test1").
					Return(nil, assert.AnError)
			},
			args: args{
				employeeID: 100,
				merchName:  "test1",
			},
			wantErr: assert.AnError,
		},
		{
			name: "error.employeeRepo.GetByIDWithLock",
			prepare: func(m *mocks) {
				m.merchRepo.EXPECT().
					GetByName(gomock.Any(), "test1").
					Return(&model.Merch{
						ID:    1,
						Name:  "test1",
						Price: 300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(100)).
					Return(nil, assert.AnError)
			},
			args: args{
				employeeID: 100,
				merchName:  "test1",
			},
			wantErr: assert.AnError,
		},
		{
			name: "error.employeeRepo.IncreaseBalance",
			prepare: func(m *mocks) {
				m.merchRepo.EXPECT().
					GetByName(gomock.Any(), "test1").
					Return(&model.Merch{
						ID:    1,
						Name:  "test1",
						Price: 300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(100)).
					Return(&model.Employee{
						ID:      100,
						Balance: 1000,
					}, nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(-300)).
					Return(assert.AnError)
			},
			args: args{
				employeeID: 100,
				merchName:  "test1",
			},
			wantErr: assert.AnError,
		},
		{
			name: "error.inventoryRepo.AddOne",
			prepare: func(m *mocks) {
				m.merchRepo.EXPECT().
					GetByName(gomock.Any(), "test1").
					Return(&model.Merch{
						ID:    1,
						Name:  "test1",
						Price: 300,
					}, nil)
				m.trManager.EXPECT().
					Do(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, do func(context.Context) error) error {
						return do(ctx)
					})
				m.employeeRepo.EXPECT().
					GetByIDWithLock(gomock.Any(), int64(100)).
					Return(&model.Employee{
						ID:      100,
						Balance: 1000,
					}, nil)
				m.employeeRepo.EXPECT().
					IncreaseBalance(gomock.Any(), int64(100), int64(-300)).
					Return(nil)
				m.inventoryRepo.EXPECT().
					AddOne(gomock.Any(), int64(100), int64(1)).
					Return(assert.AnError)
			},
			args: args{
				employeeID: 100,
				merchName:  "test1",
			},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				employeeRepo:  NewMockemployeeRepo(ctrl),
				trManager:     NewMocktrManager(ctrl),
				inventoryRepo: NewMockinventoryRepo(ctrl),
				merchRepo:     NewMockmerchRepo(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.trManager, m.employeeRepo, m.inventoryRepo, m.merchRepo)
			require.NoError(t, err)

			err = uc.Buy(context.Background(), tc.args.employeeID, tc.args.merchName)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
