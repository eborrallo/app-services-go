package storage

import (
	"app-services-go/configs"
	"database/sql"
	"fmt"
)

func MySqlConnection(config configs.DatabaseConfig) (*sql.DB, error) {
	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName)
	db, err := sql.Open("mysql", mysqlURI)
	if err != nil {
		return nil, err
	}
	return db, nil

}
