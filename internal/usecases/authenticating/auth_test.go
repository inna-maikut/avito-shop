package authenticating

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"github.com/inna-maikut/avito-shop/internal/model"
)

func TestUseCase_Auth(t *testing.T) {
	type mocks struct {
		employeeRepo  *MockemployeeRepo
		tokenProvider *MocktokenProvider
	}
	type args struct {
		username string
		password string
	}

	testCases := []struct {
		name    string
		prepare func(m *mocks)
		args    args
		wantRes string
		wantErr error
	}{
		{
			name: "success.employee_created",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(nil, model.ErrEmployeeNotFound)
				m.employeeRepo.EXPECT().
					Create(gomock.Any(), "test1", gomock.Any(), int64(1000)).
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Password: makePasswordHash("password1"),
						Balance:  1000,
					}, nil)
				m.tokenProvider.EXPECT().CreateToken("test1", int64(100)).Return("654321", nil)
			},
			args: args{
				username: "test1",
				password: "password1",
			},
			wantRes: "654321",
			wantErr: nil,
		},
		{
			name: "success.login",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Password: makePasswordHash("password1"),
						Balance:  0,
					}, nil)
				m.tokenProvider.EXPECT().CreateToken("test1", int64(100)).Return("654321", nil)
			},
			args: args{
				username: "test1",
				password: "password1",
			},
			wantRes: "654321",
			wantErr: nil,
		},
		{
			name: "error.employee_repo.get_by_username",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(nil, assert.AnError)
			},
			args: args{
				username: "test1",
				password: "password1",
			},
			wantRes: "",
			wantErr: assert.AnError,
		},
		{
			name: "error.token_provider.create_token",
			prepare: func(m *mocks) {
				m.employeeRepo.EXPECT().
					GetByUsername(gomock.Any(), "test1").
					Return(&model.Employee{
						ID:       100,
						Username: "test1",
						Password: makePasswordHash("password1"),
						Balance:  0,
					}, nil)
				m.tokenProvider.EXPECT().CreateToken("test1", int64(100)).Return("", assert.AnError)
			},
			args: args{
				username: "test1",
				password: "password1",
			},
			wantRes: "",
			wantErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := &mocks{
				employeeRepo:  NewMockemployeeRepo(ctrl),
				tokenProvider: NewMocktokenProvider(ctrl),
			}

			tc.prepare(m)

			uc, err := New(m.employeeRepo, m.tokenProvider)
			require.NoError(t, err)

			res, err := uc.Auth(context.Background(), tc.args.username, tc.args.password)
			require.ErrorIs(t, err, tc.wantErr)

			require.Equal(t, tc.wantRes, res)
		})
	}
}

func makePasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
