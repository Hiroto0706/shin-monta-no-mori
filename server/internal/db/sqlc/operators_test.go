package db_test

import (
	"context"
	db "shin-monta-no-mori/internal/db/sqlc"
	"shin-monta-no-mori/pkg/lib/password"
	"shin-monta-no-mori/pkg/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateOperator(t *testing.T) {
	hashedPassword, err := password.HashPassword("testtest")
	require.NoError(t, err)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.CreateOperatorParams
		want    db.Operator
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.CreateOperatorParams{
				Name:           "test_operator_name_10001",
				HashedPassword: hashedPassword,
				Email:          "test_10001@test.com",
			},
			want: db.Operator{
				Name:           "test_operator_name_10001",
				HashedPassword: hashedPassword,
				Email:          "test_10001@test.com",
			},
			wantErr: false,
		},
		{
			name: "異常系（nameが空文字の場合）",
			arg: db.CreateOperatorParams{
				Name: "test_operator_name_10001",
			},
			wantErr: true,
		},
		{
			name: "異常系（hashed_passwordが空文字の場合）",
			arg: db.CreateOperatorParams{
				Name:           "test_operator_name_10002",
				HashedPassword: "",
			},
			wantErr: true,
		},
		{
			name: "異常系（emailが空文字の場合）",
			arg: db.CreateOperatorParams{
				Name:           "test_operator_name_10002",
				HashedPassword: hashedPassword,
				Email:          "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operator, err := testQueries.CreateOperator(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, operator)
				require.Equal(t, tt.arg.Name, operator.Name)
				require.Equal(t, tt.arg.HashedPassword, operator.HashedPassword)
				require.Equal(t, tt.arg.Email, operator.Email)
				require.NotZero(t, operator.ID)
				require.NotZero(t, operator.CreatedAt)
			}
		})
	}
}

func TestGetOperator(t *testing.T) {
	hashedPassword, err := password.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	type args struct {
		id     int64
		params db.CreateOperatorParams
	}
	tests := []struct {
		name    string
		arg     args
		want    db.Operator
		wantErr bool
	}{
		{
			name: "正常系",
			arg: args{
				id: 10001,
				params: db.CreateOperatorParams{
					Name:           "test_operator_name_00001",
					HashedPassword: hashedPassword,
					Email:          "test_00001@test.com",
				},
			},
			want: db.Operator{
				ID:             10001,
				Name:           "test_operator_name_00001",
				HashedPassword: hashedPassword,
				Email:          "test_00001@test.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operator1, err := testQueries.CreateOperator(context.Background(), tt.arg.params)
			require.NoError(t, err)

			operator2, err := testQueries.GetOperatorByEmail(context.Background(), operator1.Email)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want.Name, operator2.Name)
				require.Equal(t, tt.want.HashedPassword, operator2.HashedPassword)
				require.Equal(t, tt.want.Email, operator2.Email)
				// IDのテストはどうしても難しい
				require.NotZero(t, operator2.ID)
				require.NotZero(t, operator2.CreatedAt)
			}
		})
	}
}

func TestUpdateOperator(t *testing.T) {
	hashedPassword, err := password.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	SetUp(t, testQueries)
	defer TearDown(t, testQueries)

	tests := []struct {
		name    string
		arg     db.UpdateOperatorParams
		want    db.Operator
		wantErr bool
	}{
		{
			name: "正常系",
			arg: db.UpdateOperatorParams{
				ID:             10001,
				Name:           "test_operator_name_10001_edited",
				HashedPassword: hashedPassword,
				Email:          "test_10001+100@test.com",
			},
			want: db.Operator{
				ID:             10001,
				Name:           "test_operator_name_10001_edited",
				HashedPassword: hashedPassword,
				Email:          "test_10001+1000@test.com",
			},
			wantErr: false,
		},
		{
			name: "異常系（存在しないIDを指定している場合）",
			arg: db.UpdateOperatorParams{
				ID: 99999,
			},
			wantErr: true,
		},
		{
			name: "異常系（titleが空になる場合）",
			arg: db.UpdateOperatorParams{
				ID:   10002,
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "異常系（hashed_passwordが空になる場合）",
			arg: db.UpdateOperatorParams{
				ID:             10003,
				HashedPassword: "",
			},
			wantErr: true,
		},
		{
			name: "異常系（emailが空になる場合）",
			arg: db.UpdateOperatorParams{
				ID:    10004,
				Email: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operator, err := testQueries.UpdateOperator(context.Background(), tt.arg)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, operator)
				require.Equal(t, tt.arg.ID, operator.ID)
				require.Equal(t, tt.arg.Name, operator.Name)
				require.Equal(t, tt.arg.HashedPassword, operator.HashedPassword)
				require.Equal(t, tt.arg.Email, operator.Email)
				require.NotZero(t, operator.CreatedAt)
			}
		})
	}
}
