package db

import (
	"context"
	"testing"
	"time"

	"github.com/haniifac/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: util.RandomString(10),
		Fullname:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordUpdatedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	expectedUser := createRandomUser(t)
	actualUser, err := testQueries.GetUser(context.Background(), expectedUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, actualUser)

	require.Equal(t, expectedUser.Username, actualUser.Username)
	require.Equal(t, expectedUser.HashedPassword, actualUser.HashedPassword)
	require.Equal(t, expectedUser.Fullname, actualUser.Fullname)
	require.Equal(t, expectedUser.Email, actualUser.Email)

	require.WithinDuration(t, expectedUser.CreatedAt, actualUser.CreatedAt, time.Second)
	require.WithinDuration(t, expectedUser.PasswordUpdatedAt, actualUser.PasswordUpdatedAt, time.Second)
}
