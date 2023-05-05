package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"vac_informer_tgbot/pkg/database/entities"
	"vac_informer_tgbot/pkg/database/telegram_db_methods"
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

var db *sql.DB

func PostgresConnDb(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Login, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSL)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return db, err
	}

	//Check for connection
	err = db.Ping()
	if err != nil {
		return db, err
	}

	telegram_db_methods.GetDatabase(db)

	return db, nil
}

func searchVacancy(vacancy *entities.VacancyInfo) bool {
	var rowHash int

	db.QueryRow("SELECT id FROM vacancy_info WHERE link_hash = $1", vacancy.LinkHash).
		Scan(&rowHash)
	if rowHash != 0 {
		return true
	}

	return false
}

func createVacancy(vacancy *entities.VacancyInfo) {
	query := "INSERT INTO vacancy_info (website, vacancy_link, vacancy_text, link_hash, location, company) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(query, vacancy.Website,
		vacancy.VacancyLink, vacancy.VacancyText, vacancy.LinkHash, vacancy.Location, vacancy.Company)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckVacancy(website, vacancyLink, vacancyTitle, location, company, hash, text string) {
	vacancy := &entities.VacancyInfo{
		Website:     website,
		VacancyLink: vacancyLink,
		VacancyText: vacancyTitle,
		Location:    location,
		Company:     company,
		LinkHash:    hash,
	}

	searchVac := searchVacancy(vacancy)
	if searchVac == false {
		createVacancy(vacancy)
		telegram.SendMessage(text)
	}
}
