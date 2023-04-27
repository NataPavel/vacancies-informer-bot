package database

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
	"vac_informer_tgbot/pkg/database/entities"
	"vac_informer_tgbot/pkg/telegram"
)

type Config struct {
	Login    string
	Password string
	DBName   string
	Host     string
	Port     string
	SSL      string
}

func PostgresConnDb(cfg Config) (*sql.DB, error) {
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

func searchVacancy(db *sql.DB, vacancy *entities.VacancyInfo) bool {
	var rowHash int

	db.QueryRow("SELECT id FROM vacancy_info WHERE link_hash = $1", vacancy.LinkHash).
		Scan(&rowHash)
	if rowHash != 0 {
		return true
	}

	return false
}

func createVacancy(db *sql.DB, vacancy *entities.VacancyInfo) {
	query := "INSERT INTO vacancy_info (website, vacancy_link, vacancy_text, link_hash, location, company) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(query, vacancy.Website,
		vacancy.VacancyLink, vacancy.VacancyText, vacancy.LinkHash, vacancy.Location, vacancy.Company)
	if err != nil {
		log.Fatal(err)
	}
}

func ErrorLogger(errorMessage string, db *sql.DB) {
	errorLog := &entities.ErrorLogger{
		ErrorMessage: errorMessage,
	}

	query := "INSERT INTO error_log (error_message) VALUES ($1)"

	_, err := db.Exec(query, errorLog.ErrorMessage)
	if err != nil {
		log.Fatalf("Failed to register an error: %s", err)
	}
}

func CheckVacancy(website, vacancyLink, vacancyTitle, location, company, hash, text string, db *sql.DB, tgbot *tgbotapi.BotAPI) {
	vacancy := &entities.VacancyInfo{
		Website:     website,
		VacancyLink: vacancyLink,
		VacancyText: vacancyTitle,
		Location:    location,
		Company:     company,
		LinkHash:    hash,
	}

	searchVac := searchVacancy(db, vacancy)
	if searchVac == false {
		createVacancy(db, vacancy)
		telegram.SendMessage(text, tgbot)
	}

}
