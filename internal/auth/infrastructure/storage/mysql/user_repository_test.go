package mysql

import (
	"app-services-go/configs"
	user "app-services-go/internal/auth/domain"
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UserRepository_Save_RepositoryError(t *testing.T) {
	userID, userName, userEmail, userPassword := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test User", "aaa@gmail.com", "123123"
	user, err := user.NewUser(userID, userName, userEmail, userPassword)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO users (id, name, email, password, validated) VALUES (?, ?, ?, ?, ?)").
		WithArgs(userID, userName, userEmail, userPassword, false).
		WillReturnError(errors.New("something-failed"))
	c, _ := configs.GetDatabaseConfig()
	repo := NewUserRepository(db, c)

	err = repo.Save(context.Background(), user)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.Error(t, err)
}

func Test_UserRepository_Save_Succeed(t *testing.T) {
	userID, userName, userEmail, userPassword := "37a0f027-15e6-47cc-a5d2-64183281087e", "Test User", "aaa@gmail.com", "123123"

	user, err := user.NewUser(userID, userName, userEmail, userPassword)
	require.NoError(t, err)

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlMock.ExpectExec(
		"INSERT INTO users (id, name, email, password, validated) VALUES (?, ?, ?, ?, ?)").
		WithArgs(userID, userName, userEmail, userPassword, false).
		WillReturnResult(sqlmock.NewResult(0, 1))
	c, _ := configs.GetDatabaseConfig()

	repo := NewUserRepository(db, c)

	err = repo.Save(context.Background(), user)

	assert.NoError(t, sqlMock.ExpectationsWereMet())
	assert.NoError(t, err)
}
