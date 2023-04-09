package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"vac_informer_tgbot/pkg/database/entities"
	"vac_informer_tgbot/pkg/telegram"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	Login    string
	Password string
	DBName   string
	Host     string
	Port     string
	SSL      string
}

var vacancy *entities.VacancyInfo

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

func searchVacancy(db *sql.DB) bool {
	var rowHash int

	db.QueryRow("SELECT id FROM vacancy_info WHERE link_hash = $1", vacancy.LinkHash).
		Scan(&rowHash)
	if rowHash != 0 {
		return true
	}

	return false
}

func createVacancy(db *sql.DB) {
	query := "INSERT INTO vacancy_info (website, vacancy_link, vacancy_text, link_hash) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(query, vacancy.Website,
		vacancy.VacancyLink, vacancy.VacancyText, vacancy.LinkHash)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckVacancy(website, vacancyLink, vacancyText, hash string) {
	db, err := connDb(Config{})
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	vacancy.Website = website
	vacancy.VacancyLink = vacancyLink
	vacancy.VacancyText = vacancyText
	vacancy.LinkHash = hash

	searchVac := searchVacancy(db)
	if searchVac == false {
		createVacancy(db)
		telegram.SendMessage(vacancyText)
	}
}
