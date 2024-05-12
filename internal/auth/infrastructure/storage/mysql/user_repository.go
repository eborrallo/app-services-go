package mysql

import (
	"app-services-go/configs"
	user "app-services-go/internal/auth/domain"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
)

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

// Save implements the user.UserRepository interface.
func (r *UserRepository) Save(ctx context.Context, user user.User) error {
	userSQLStruct := sqlbuilder.NewStruct(new(sqlUser))
	query, args := userSQLStruct.InsertInto(sqlUserTable, sqlUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist user on database: %v", err)
	}

	return nil
}

// Update implements the user.UserRepository interface.
func (r *UserRepository) Update(ctx context.Context, user user.User) error {
	userSQLStruct := sqlbuilder.NewStruct(new(sqlUser))
	query, args := userSQLStruct.Update(sqlUserTable, sqlUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Validated: user.Validated,
	}).Where(fmt.Sprintf("ID = %s", user.ID)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(ctxTimeout, query, args...)
	if err != nil {
		return fmt.Errorf("error trying to persist user on database: %v", err)
	}

	return nil
}

func (r *UserRepository) FetchByEmail(ctx context.Context, email string) (user.User, error) {

	var entity user.User

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("users")
	sb.Where(sb.Equal("Email", email))

	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctxTimeout, query, args...).Scan(&entity.ID, &entity.Name, &entity.Email, &entity.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("error trying to fetch user %s on database: %v", email, err)
	}

	return entity, nil
}

func (r *UserRepository) FetchById(ctx context.Context, id string) (user.User, error) {

	var entity user.User

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("users")
	sb.Where(sb.Equal("Id", id))

	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctxTimeout, query, args...).Scan(&entity.ID, &entity.Name, &entity.Email, &entity.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("error trying to fetch user %s on database: %v", id, err)
	}

	return entity, nil
}

func (r *UserRepository) FetchByAddress(ctx context.Context, address string) (user.User, error) {

	var entity user.User

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("*")
	sb.From("users")
	sb.Where(sb.Equal("Address", address))

	query, args := sb.Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.QueryRowContext(ctxTimeout, query, args...).Scan(&entity.ID, &entity.Name, &entity.Email, &entity.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("error trying to fetch user %s on database: %v", address, err)
	}

	return entity, nil
}
