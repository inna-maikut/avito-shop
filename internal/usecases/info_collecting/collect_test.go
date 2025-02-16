package info_collecting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/inna-maikut/avito-shop/internal/model"
)

func TestUseCase_Collect(t *testing.T) {
	type mocks struct {
		employeeRepo    *MockemployeeRepo
		transactionRepo *MocktransactionRepo
		inventoryRepo   *MockinventoryRepo
	}
	type args struct {
		employeeID int64
	}

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantRes model.EmployeeInfo
		wantErr error
	}{
		{
			name: "success.collection",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByID(gomock.Any(), int64(100)).
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  1000,
					}, nil).AnyTimes()
				m.transactionRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return([]model.Transaction{
						{
							IsSender:               true,
							CounterpartyEmployeeID: 200,
							CounterpartyUsername:   "test2",
							Amount:                 100,
						},
						{
							IsSender:               false,
							CounterpartyEmployeeID: 300,
							CounterpartyUsername:   "test3",
							Amount:                 100,
						},
						{
							IsSender:               true,
							CounterpartyEmployeeID: 400,
							CounterpartyUsername:   "test4",
							Amount:                 100,
						},
					}, nil).AnyTimes()
				m.inventoryRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return([]model.Inventory{
						{
							EmployeeID: 100,
							MerchID:    1,
							Quantity:   1,
							MerchName:  "cup",
						},
						{
							EmployeeID: 100,
							MerchID:    2,
							Quantity:   1,
							MerchName:  "book",
						},
						{
							EmployeeID: 100,
							MerchID:    3,
							Quantity:   10,
							MerchName:  "socks",
						},
					}, nil).AnyTimes()
			},
			args: args{
				employeeID: 100,
			},
			wantRes: model.EmployeeInfo{
				Coins: 1000,
				Inventory: []model.Inventory{
					{
						EmployeeID: 100,
						MerchID:    1,
						Quantity:   1,
						MerchName:  "cup",
					},
					{
						EmployeeID: 100,
						MerchID:    2,
						Quantity:   1,
						MerchName:  "book",
					},
					{
						EmployeeID: 100,
						MerchID:    3,
						Quantity:   10,
						MerchName:  "socks",
					},
				},
				ReceivedTransactions: []model.Transaction{
					{
						IsSender:               false,
						CounterpartyEmployeeID: 300,
						CounterpartyUsername:   "test3",
						Amount:                 100,
					},
				},
				SentTransactions: []model.Transaction{
					{
						IsSender:               true,
						CounterpartyEmployeeID: 200,
						CounterpartyUsername:   "test2",
						Amount:                 100,
					},
					{
						IsSender:               true,
						CounterpartyEmployeeID: 400,
						CounterpartyUsername:   "test4",
						Amount:                 100,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "success.collection.zeroinformation",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByID(gomock.Any(), int64(100)).
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  1000,
					}, nil).AnyTimes()
				m.transactionRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return([]model.Transaction{}, nil).AnyTimes()
				m.inventoryRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return([]model.Inventory{}, nil).AnyTimes()
			},
			args: args{
				employeeID: 100,
			},
			wantRes: model.EmployeeInfo{
				Coins:                1000,
				Inventory:            []model.Inventory{},
				ReceivedTransactions: []model.Transaction{},
				SentTransactions:     []model.Transaction{},
			},
			wantErr: nil,
		},
		{
			name: "error.inventoryRepo.GetByEmployee",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByID(gomock.Any(), int64(100)).
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  1000,
					}, nil).AnyTimes()
				m.transactionRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return([]model.Transaction{
						{
							IsSender:               true,
							CounterpartyEmployeeID: 200,
							CounterpartyUsername:   "test2",
							Amount:                 100,
						},
						{
							IsSender:               false,
							CounterpartyEmployeeID: 300,
							CounterpartyUsername:   "test3",
							Amount:                 100,
						},
						{
							IsSender:               true,
							CounterpartyEmployeeID: 400,
							CounterpartyUsername:   "test4",
							Amount:                 100,
						},
					}, nil).AnyTimes()
				m.inventoryRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return(nil, assert.AnError).AnyTimes()
			},
			args: args{
				employeeID: 100,
			},
			wantRes: model.EmployeeInfo{},
			wantErr: assert.AnError,
		},
		{
			name: "error.transactionRepo.GetByEmployee",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByID(gomock.Any(), int64(100)).
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Balance:  1000,
					}, nil).AnyTimes()
				m.transactionRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return(nil, assert.AnError).AnyTimes()
				m.inventoryRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return([]model.Inventory{
						{
							EmployeeID: 100,
							MerchID:    1,
							Quantity:   1,
							MerchName:  "cup",
						},
						{
							EmployeeID: 100,
							MerchID:    2,
							Quantity:   1,
							MerchName:  "book",
						},
						{
							EmployeeID: 100,
							MerchID:    3,
							Quantity:   10,
							MerchName:  "socks",
						},
					}, nil).AnyTimes()
			},
			args: args{
				employeeID: 100,
			},
			wantRes: model.EmployeeInfo{},
			wantErr: assert.AnError,
		},
		{
			name: "error.employeeRepo.GetByID",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByID(gomock.Any(), int64(100)).
					Return(nil, assert.AnError).AnyTimes()
				m.transactionRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return([]model.Transaction{
						{
							IsSender:               true,
							CounterpartyEmployeeID: 200,
							CounterpartyUsername:   "test2",
							Amount:                 100,
						},
						{
							IsSender:               false,
							CounterpartyEmployeeID: 300,
							CounterpartyUsername:   "test3",
							Amount:                 100,
						},
						{
							IsSender:               true,
							CounterpartyEmployeeID: 400,
							CounterpartyUsername:   "test4",
							Amount:                 100,
						},
					}, nil).AnyTimes()
				m.inventoryRepo.EXPECT().
					GetByEmployee(gomock.Any(), int64(100)).
					Return([]model.Inventory{
						{
							EmployeeID: 100,
							MerchID:    1,
							Quantity:   1,
							MerchName:  "cup",
						},
						{
							EmployeeID: 100,
							MerchID:    2,
							Quantity:   1,
							MerchName:  "book",
						},
						{
							EmployeeID: 100,
							MerchID:    3,
							Quantity:   10,
							MerchName:  "socks",
						},
					}, nil).AnyTimes()
			},
			args: args{
				employeeID: 100,
			},
			wantRes: model.EmployeeInfo{},
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				employeeRepo:    NewMockemployeeRepo(ctrl),
				transactionRepo: NewMocktransactionRepo(ctrl),
				inventoryRepo:   NewMockinventoryRepo(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.employeeRepo, m.transactionRepo, m.inventoryRepo)
			require.NoError(t, err)

			res, err := uc.Collect(context.Background(), tc.args.employeeID)

			require.ErrorIs(t, err, tc.wantErr)

			require.Equal(t, tc.wantRes, res)
		})
	}
}
