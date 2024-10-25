package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/yourusername/tamagotchi-bot/internal/bot"
	"github.com/yourusername/tamagotchi-bot/internal/pet"
)

func main() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	// Инициализация базы данных
	pet.InitDB("pets.db")

	// Запуск веб-сервера в отдельной горутине
	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/webapp/data", handleWebAppData).Methods("POST")
		router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

		log.Println("Запуск веб-сервера на порту :8080")
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			log.Fatal("Ошибка запуска веб-сервера: ", err)
		}
	}()

	// Запуск Telegram бота
	bot.Start()
}

// Обработчик POST-запроса от WebApp для получения данных о питомце
func handleWebAppData(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID    int64 `json:"user_id"`
		Hunger    int   `json:"hunger"`
		Happiness int   `json:"happiness"`
		Coins     int   `json:"coins"`
	}

	// Декодируем JSON запрос
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создаем объект питомца с данными из запроса
	petState := &pet.Pet{
		Hunger:    requestData.Hunger,
		Happiness: requestData.Happiness,
		Coins:     requestData.Coins,
	}

	// Сохраняем состояние питомца в базе данных
	err = pet.SavePetState(requestData.UserID, petState)
	if err != nil {
		http.Error(w, "Ошибка сохранения состояния питомца", http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]string{"status": "success"}
	json.NewEncoder(w).Encode(response)
}
