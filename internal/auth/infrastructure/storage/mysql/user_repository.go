package mysql

import (
	"app-services-go/configs"
	user "app-services-go/internal/auth/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

// userColumns is the explicit column list, matching the Scan order below. Never
// use SELECT * with a partial Scan: it breaks the moment the table changes.
var userColumns = []string{"id", "name", "email", "password", "validated", "address"}

// UserRepository is a MySQL user.UserRepository implementation.
type UserRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
}

// NewUserRepository initializes a MySQL-based implementation of user.UserRepository.
func NewUserRepository(db *sql.DB, configs configs.DatabaseConfig) *UserRepository {
	return &UserRepository{
		db:        db,
		dbTimeout: configs.DbTimeout,
	}
}

func toSQLUser(u user.User) sqlUser {
	return sqlUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		Validated: u.Validated,
		Address:   u.Address,
	}
}

// Save implements the user.UserRepository interface.
func (r *UserRepository) Save(ctx context.Context, u user.User) error {
	userSQLStruct := sqlbuilder.NewStruct(new(sqlUser))
	query, args := userSQLStruct.InsertInto(sqlUserTable, toSQLUser(u)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(ctxTimeout, query, args...); err != nil {
		return fmt.Errorf("error trying to persist user on database: %w", err)
	}

	return nil
}

// Update implements the user.UserRepository interface.
func (r *UserRepository) Update(ctx context.Context, u user.User) error {
	userSQLStruct := sqlbuilder.NewStruct(new(sqlUser))
	ub := userSQLStruct.Update(sqlUserTable, toSQLUser(u))
	// Parametrized predicate: no more fmt.Sprintf into the SQL string.
	ub.Where(ub.Equal("id", u.ID))
	query, args := ub.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(ctxTimeout, query, args...); err != nil {
		return fmt.Errorf("error trying to update user on database: %w", err)
	}

	return nil
}

func (r *UserRepository) fetchBy(ctx context.Context, column, value string) (user.User, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(userColumns...)
	sb.From(sqlUserTable)
	sb.Where(sb.Equal(column, value))

	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var entity user.User
	err := r.db.QueryRowContext(ctxTimeout, query, args...).
		Scan(&entity.ID, &entity.Name, &entity.Email, &entity.Password, &entity.Validated, &entity.Address)
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, user.ErrUserNotFound
	}
	if err != nil {
		return user.User{}, fmt.Errorf("error trying to fetch user %s on database: %w", value, err)
	}

	return entity, nil
}

func (r *UserRepository) FetchByEmail(ctx context.Context, email string) (user.User, error) {
	return r.fetchBy(ctx, "email", email)
}

func (r *UserRepository) FetchById(ctx context.Context, id string) (user.User, error) {
	return r.fetchBy(ctx, "id", id)
}

func (r *UserRepository) FetchByAddress(ctx context.Context, address string) (user.User, error) {
	return r.fetchBy(ctx, "address", address)
}
