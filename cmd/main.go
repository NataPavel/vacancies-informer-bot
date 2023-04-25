package main

import (
	"log"
	"os"

	"vac_informer_tgbot/pkg/database"
	"vac_informer_tgbot/pkg/services"
	"vac_informer_tgbot/pkg/telegram"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func main() {
	println("It's working")
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env variables: %s", err.Error())
	}

	db, err := database.PostgresConnDb(database.Config{
		Login:    viper.GetString("db.login"),
		Password: os.Getenv("DB_PASS"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		SSL:      viper.GetString("db.sslMode"),
		DBName:   viper.GetString("db.dbName"),
	})
	if err != nil {
		log.Fatalf("Failed connection to DB: %s", err)
	}
	defer db.Close()

	tgbot, err := telegram.NewTgBot(telegram.Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		TelegramUrl: viper.GetString("telegramURL"),
	})

	c := cron.New()
	c.AddFunc("*/5 * * * *", func() {
		searchTags := []string{"Golang", ".Net"}
		for _, i := range searchTags {
			services.Indeed(i, db, tgbot)
			services.Dou(i, db, tgbot)
			services.Jooble(i, db, tgbot)
			services.Djinni(i, db, tgbot)
		}
	})
	c.Start()

	select {}

}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
