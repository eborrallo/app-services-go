package storage

import (
	"app-services-go/configs"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MySqlURI(config configs.DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName)
}

func MySqlConnection(config configs.DatabaseConfig) (*sql.DB, error) {
	mysqlURI := MySqlURI(config)
	db, err := sql.Open("mysql", mysqlURI)
	if err != nil {
		return nil, err
	}

	// sql.Open does not actually connect; configure the pool and ping so a bad
	// connection fails fast at startup instead of on the first query.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("cannot reach database: %w", err)
	}

	// Migrations are advisory here (the schema is also seeded via the
	// docker-compose init-dir mount). Surface failures instead of swallowing
	// them, but do not abort startup over them.
	if err := doMigrations(config); err != nil {
		log.Printf("warning: migrations not applied: %v", err)
	}

	return db, nil
}

func doMigrations(config configs.DatabaseConfig) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current working directory: %w", err)
	}
	migrationsDir := filepath.Join(currentDir, "sql")

	m, err := migrate.New("file://"+migrationsDir, "mysql://"+MySqlURI(config))
	if err != nil {
		return fmt.Errorf("initializing migrations: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}
