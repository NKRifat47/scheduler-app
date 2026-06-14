package db

import (
	"fmt"
	"scheduler-app/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetConnectionString(cnf *config.DBConfig) string {
	sslmode := "disable"
	if cnf.EnableSSLMode {
		sslmode = "require"
	}

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.Name, sslmode,
	)
}

func NewConnection(cnf *config.DBConfig) (*sqlx.DB, error) {
	dbSource := GetConnectionString(cnf)
	dbCon, err := sqlx.Connect("postgres", dbSource)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return dbCon, nil
}