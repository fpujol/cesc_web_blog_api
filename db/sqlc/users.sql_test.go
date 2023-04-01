package db

import (
	"blogapi/pkg/password"
	"blogapi/pkg/utils"
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	email := utils.RandomEmail()
	hashedPassword, err := password.HashPassword(email, utils.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Email:          utils.RandomEmail(),
		FirstName:      utils.RandomString(6),
		LastName:       utils.RandomString(6),
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
		CreatedBy:      "Francesc Pujol",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByEmail(context.Background(), user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}


func TestUpdateUserByEmail(t *testing.T) {
	type args struct {
		ctx context.Context
		arg UpdateUserByEmailParams
	}
	tests := []struct {
		name    string
		q       *Queries
		args    args
		want    User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.UpdateUserByEmail(tt.args.ctx, tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Queries.UpdateUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Queries.UpdateUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
