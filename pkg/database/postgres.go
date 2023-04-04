package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Login    string
	Password string
	DBName   string
	Host     string
	Port     string
	SSL      string
}

func connDb(cfg Config) (*sql.DB, error) {
	cfg.Login = viper.GetString("db.login")
	cfg.Password = os.Getenv("DB_PASS")
	cfg.Host = viper.GetString("db.host")
	cfg.Port = viper.GetString("db.port")
	cfg.SSL = viper.GetString("db.sslMode")
	cfg.DBName = viper.GetString("db.dbName")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Login, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSL)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	//Check for connection
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
