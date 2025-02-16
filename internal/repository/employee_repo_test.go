//go:build integration

package repository

import (
	"context"
	"testing"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/stretchr/testify/require"

	"github.com/inna-maikut/avito-shop/internal/model"
)

func Test_GetByUsername(t *testing.T) {
	db := setUp(t)
	repo, err := NewEmployeeRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)

	type args struct {
		username string
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		wantRes *model.Employee
		wantErr error
	}{
		{
			name: "found",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM employee where username = $1`, "get-by-username-1")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO employee (username, password, balance)
					VALUES ($1, $2, $3)`, "get-by-username-1", "password", 500)
				require.NoError(t, err)
			},
			args: args{
				username: "get-by-username-1",
			},
			wantRes: &model.Employee{
				Username: "get-by-username-1",
				Password: "password",
				Balance:  500,
			},
			wantErr: nil,
		},
		{
			name: "not_found",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM employee where username = $1`, "get-by-username-2")
				require.NoError(t, err)
			},
			args: args{
				username: "get-by-username-2",
			},
			wantRes: nil,
			wantErr: model.ErrEmployeeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.GetByUsername(context.Background(), tc.args.username)

			require.ErrorIs(t, err, tc.wantErr)

			if res != nil {
				res.ID = 0 // can't validate id, because it's autoincrement
			}
			require.Equal(t, tc.wantRes, res)
		})
	}
}
