package storage

import (
	"app-services-go/configs"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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

	doMigrations(config)

	return db, nil

}

func doMigrations(config configs.DatabaseConfig) error {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return err
	}
	migrationsDir := filepath.Join(currentDir, "sql")

	m, err := migrate.New("file://"+migrationsDir, "mysql://"+MySqlURI(config))

	if err != nil {
		fmt.Println("Error running migrations:", err)
		return err
	}
	m.Up()

	return nil
}
