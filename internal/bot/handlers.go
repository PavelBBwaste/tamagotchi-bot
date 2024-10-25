package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleUpdate обрабатывает обновления от Telegram
func handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		handleCommand(update.Message)
	} else if update.CallbackQuery != nil {
		handleCallbackQuery(update.CallbackQuery)
	}
}

// handleCommand обрабатывает команды от пользователей
func handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		sendStartMessage(message.Chat.ID)

	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда.")
		Bot.Send(msg)
	}
}

// sendStartMessage отправляет приветственное сообщение с кнопкой для открытия игры
func sendStartMessage(chatID int64) {
	// Создание кнопки для веб-приложения
	webAppInfo := tgbotapi.WebAppInfo{URL: "https://pavelbbwaste.github.io/testsite/"}
	webAppButton := tgbotapi.NewInlineKeyboardButtonWebApp("Открыть игру", webAppInfo)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(webAppButton),
	)

	msg := tgbotapi.NewMessage(chatID, "Добро пожаловать в Тамагочи!")
	msg.ReplyMarkup = keyboard

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}

// handleCallbackQuery обрабатывает запросы обратного вызова от кнопок
func handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {

	// Дополнительно: отправляем сообщение в чат
	msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Игра открыта!")
	if _, err := Bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
}
