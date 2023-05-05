package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vac_informer_tgbot/pkg/database"
	"vac_informer_tgbot/pkg/services"
	"vac_informer_tgbot/pkg/telegram"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

func main() {
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
			services.Indeed(i)
			services.Dou(i)
			services.Jooble(i)
			services.Djinni(i)
		}
	})
	c.Start()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-shutdownChan
		log.Println("Received shutdown signal. Shutting down server...")

		// Set a timeout of 5 seconds to allow the server to shut down gracefully
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = c.Stop().Err(); err != nil {
			log.Printf("Failed to stop cron: %s", err)
		}

		if err = db.Close(); err != nil {
			log.Fatalf("Failed to close db connection: %s", err)
		}

		tgbot.StopReceivingUpdates()

		log.Println("Server has shut down gracefully")
		os.Exit(0)
	}()

	log.Println("Server is up and running")

	select {}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
