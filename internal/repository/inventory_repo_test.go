//go:build integration

package repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-shop/internal/model"
)

func Test_GetByEmployee(t *testing.T) {
	db := setUp(t)
	repo, err := NewInventoryRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)

	type args struct {
		employeeID int64
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		wantRes []model.Inventory
		wantErr error
	}{
		{
			name: "get",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM inventory where employee_id = $1`, 390294)
				require.NoError(t, err)
				_, err = db.Exec(`DELETE FROM inventory where employee_id = $1`, 390295)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO inventory (employee_id, merch_id, quantity)
					VALUES ($1, $2, 1)`, 390294, 1)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO inventory (employee_id, merch_id, quantity)
					VALUES ($1, $2, 1)`, 390294, 2)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO inventory (employee_id, merch_id, quantity)
					VALUES ($1, $2, 1)`, 390295, 2)
				require.NoError(t, err)
			},
			args: args{
				employeeID: 390294,
			},
			wantRes: []model.Inventory{
				{
					EmployeeID: 390294,
					MerchID:    1,
					Quantity:   1,
					MerchName:  "t-shirt",
				},
				{
					EmployeeID: 390294,
					MerchID:    2,
					Quantity:   1,
					MerchName:  "cup",
				},
			},
			wantErr: nil,
		},
		{
			name:    "empty",
			prepare: func(_ *testing.T) {},
			args: args{
				employeeID: 390288,
			},
			wantRes: []model.Inventory{},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.GetByEmployee(context.Background(), tc.args.employeeID)

			require.ErrorIs(t, err, tc.wantErr)

			require.Equal(t, tc.wantRes, res)
		})
	}
}

func Test_AddOne(t *testing.T) {
	db := setUp(t)
	repo, err := NewInventoryRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)

	type args struct {
		employeeID int64
		merchID    int64
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		check   func(t *testing.T)
		wantErr error
	}{
		{
			name: "get",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM inventory where employee_id = $1`, 390294)
				require.NoError(t, err)
				_, err = db.Exec(`DELETE FROM inventory where employee_id = $1`, 390295)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO inventory (employee_id, merch_id, quantity)
					VALUES ($1, $2, 1)`, 390294, 1)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO inventory (employee_id, merch_id, quantity)
					VALUES ($1, $2, 1)`, 390294, 2)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO inventory (employee_id, merch_id, quantity)
					VALUES ($1, $2, 1)`, 390295, 2)
				require.NoError(t, err)
			},
			args: args{
				employeeID: 390294,
				merchID:    1,
			},
			check: func(t *testing.T) {
				var inventory InventoryWithMerchName
				err = db.Get(&inventory, "SELECT employee_id, merch_id, quantity, merch_name FROM inventory WHERE employee_id = $1 AND merch_id = $2", 390294, 1)
				require.NoError(t, err)

				require.Equal(t, InventoryWithMerchName{
					EmployeeID: 390294,
					MerchID:    1,
					Quantity:   2,
					MerchName:  "t-shirt",
				}, inventory)
			},
			wantErr: nil,
		},
		{
			name:    "empty",
			prepare: func(_ *testing.T) {},
			args: args{
				employeeID: 390288,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			err := repo.AddOne(context.Background(), tc.args.employeeID, tc.args.merchID)

			require.ErrorIs(t, err, tc.wantErr)

			// require.Equal(t, tc.wantRes, res)
		})
	}
}
