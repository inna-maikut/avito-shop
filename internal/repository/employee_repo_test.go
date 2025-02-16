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

func Test_GetByID(t *testing.T) {
	db := setUp(t)
	repo, err := NewEmployeeRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)

	const ID = 46875401

	type args struct {
		employeeID int64
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
				_, err = db.Exec(`DELETE FROM employee where username = $1`, "get-by-username-100")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO employee (id, username, password, balance)
					VALUES ($1, $2, $3, $4)`, ID, "get-by-username-100", "password", 500)
				require.NoError(t, err)
			},
			args: args{
				employeeID: ID,
			},
			wantRes: &model.Employee{
				ID:       ID,
				Username: "get-by-username-100",
				Password: "password",
				Balance:  500,
			},
			wantErr: nil,
		},
		{
			name: "not_found",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM employee where id = $1`, ID+1)
				require.NoError(t, err)
			},
			args: args{
				employeeID: ID + 1,
			},
			wantRes: nil,
			wantErr: model.ErrEmployeeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.GetByID(context.Background(), tc.args.employeeID)

			require.ErrorIs(t, err, tc.wantErr)

			require.Equal(t, tc.wantRes, res)
		})
	}
}

func Test_Create(t *testing.T) {
	db := setUp(t)
	repo, err := NewEmployeeRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)

	const (
		ID       = 46825601
		username = "create-46825601"
	)

	type args struct {
		username     string
		passwordHash string
		balance      int64
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
				_, err = db.Exec(`DELETE FROM employee where username = $1`, "get-by-username-3")
				require.NoError(t, err)

			},
			args: args{
				username:     "get-by-username-3",
				passwordHash: "password",
				balance:      500,
			},
			wantRes: &model.Employee{
				Username: "get-by-username-3",
				Password: "password",
				Balance:  500,
			},
			wantErr: nil,
		},
		{
			name: "AlreadyExists",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM employee where id = $1`, ID)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO employee (id, username, password, balance)
					VALUES ($1, $2, $3, $4)`, ID, username, "password", 500)
				require.NoError(t, err)
			},
			args: args{
				username:     username,
				passwordHash: "password",
				balance:      500,
			},
			wantRes: nil,
			wantErr: model.ErrEmployeeAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.Create(context.Background(), tc.args.username, tc.args.passwordHash, tc.args.balance)

			require.ErrorIs(t, err, tc.wantErr)

			if res != nil {
				res.ID = 0 // can't validate id, because it's autoincrement
			}

			require.Equal(t, tc.wantRes, res)
		})
	}
}

func Test_GetByIDWithLock(t *testing.T) {
	db := setUp(t)
	repo, err := NewEmployeeRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)

	const ID = 46863401

	type args struct {
		employeeID int64
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
				_, err = db.Exec(`DELETE FROM employee where username = $1`, "get-by-username-5")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO employee (id, username, password, balance)
					VALUES ($1, $2, $3, $4)`, ID, "get-by-username-5", "password", 500)
				require.NoError(t, err)
			},
			args: args{
				employeeID: ID,
			},
			wantRes: &model.Employee{
				ID:       ID,
				Username: "get-by-username-5",
				Password: "password",
				Balance:  500,
			},
			wantErr: nil,
		},
		{
			name: "not_found",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM employee where id = $1`, ID+1)
				require.NoError(t, err)
			},
			args: args{
				employeeID: ID + 1,
			},
			wantRes: nil,
			wantErr: model.ErrEmployeeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			res, err := repo.GetByIDWithLock(context.Background(), tc.args.employeeID)

			require.ErrorIs(t, err, tc.wantErr)

			require.Equal(t, tc.wantRes, res)
		})
	}
}
func Test_IncreaseBalance(t *testing.T) {
	db := setUp(t)
	repo, err := NewEmployeeRepository(db, trmsqlx.DefaultCtxGetter)
	require.NoError(t, err)

	const ID = 46825600

	type args struct {
		employeeID int64
		amount     int64
	}

	testCases := []struct {
		name    string
		prepare func(t *testing.T)
		args    args
		check   func(t *testing.T)
		wantErr error
	}{
		{
			name: "found",
			prepare: func(t *testing.T) {
				_, err = db.Exec(`DELETE FROM employee where username = $1`, "get-by-username-6")
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO employee (id, username, password, balance)
					VALUES ($1, $2, $3, $4)`, ID, "get-by-username-6", "password", 1000)
				require.NoError(t, err)

			},
			args: args{
				employeeID: ID,
				amount:     500,
			},
			check: func(t *testing.T) {
				var employee Employee
				err = db.Get(&employee, "SELECT id, username, password, balance FROM employee WHERE username = $1", "get-by-username-6")
				require.NoError(t, err)

				require.Equal(t, Employee{
					ID:       ID,
					Username: "get-by-username-6",
					Password: "password",
					Balance:  1500,
				}, employee)
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(t)

			err := repo.IncreaseBalance(context.Background(), tc.args.employeeID, tc.args.amount)

			require.ErrorIs(t, err, tc.wantErr)

			tc.check(t)
		})
	}
}
