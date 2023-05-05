package telegram_db_methods

import (
	"database/sql"
	"log"
	"vac_informer_tgbot/pkg/database/entities"
)

var db *sql.DB

func GetDatabase(database *sql.DB) {
	db = database
}

func ErrorLogger(errorMessage string) {
	errorLog := &entities.ErrorLogger{
		ErrorMessage: errorMessage,
	}

	query := "INSERT INTO error_log (error_message) VALUES ($1)"

	_, err := db.Exec(query, errorLog.ErrorMessage)
	if err != nil {
		log.Fatalf("Failed to register an error: %s", err)
	}
}

func searchUserChat(userChatId *entities.UserChatId) bool {
	var chatId int64

	db.QueryRow("SELECT id FROM user_chat_id WHERE chat_id = $1", userChatId.ChatId).
		Scan(&chatId)
	if chatId != 0 {
		return true
	}

	return false
}

func createUser(userChatId *entities.UserChatId) {
	query := "INSERT INTO user_chat_id (chat_id) VALUES ($1)"
	_, err := db.Exec(query, userChatId.ChatId)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckUser(chatid int64) {
	userChatId := &entities.UserChatId{
		ChatId: chatid,
	}

	searchUser := searchUserChat(userChatId)
	if searchUser == false {
		createUser(userChatId)
	}
}

func SelectAllUsers() ([]int64, error) {
	var chatIds []int64

	rows, err := db.Query("SELECT chat_id FROM user_chat_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chatId int64
		err = rows.Scan(&chatId)
		if err != nil {
			return nil, err
		}
		chatIds = append(chatIds, chatId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chatIds, nil
}
