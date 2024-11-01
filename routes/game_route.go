package routes

import (
	"encoding/json"
	"fmt"
	"golang/services"
	"golang/structs"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func LastByCode(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var req structs.LastByCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request: Verify any missing fields", http.StatusBadRequest)
		return
	}

	var game structs.Game
	err := db.Where("game_code = ?", req.Code).
		Order("created_at DESC").
		First(&game).Error

	if err == gorm.ErrRecordNotFound {
		http.Error(w, "No game found with the provided code", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve game from the database", http.StatusInternalServerError)
		return
	}

	var gameData structs.GameData
	err = db.Where("game_id = ?", game.GameID).
		Order("game_time DESC").
		First(&gameData).Error

	if err == gorm.ErrRecordNotFound {
		http.Error(w, "No game data found with the provided code", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve game data from the database", http.StatusInternalServerError)
		return
	}

	var gameProps *structs.GameProps
	if err := json.Unmarshal([]byte(gameData.GameData), &gameProps); err != nil {
		http.Error(w, "Failed to deserialize game data", http.StatusInternalServerError)
		return
	}

	start := time.Now()

	calculatedGame := services.Calculate(gameProps, req.Item)

	elapsed := time.Since(start)
	fmt.Printf("It took %s\n", elapsed)

	calculatedGameJSON, err := json.Marshal(calculatedGame)
	if err != nil {
		http.Error(w, "Failed to serialize game data", http.StatusInternalServerError)
		return
	}

	response := structs.LastByCodeResponse{
		Success: true,
		Data: structs.LastByCodeResponseData{
			GameID:       game.GameID,
			SummonerName: *game.SummonerName,
			CreatedAt:    game.CreatedAt,
			GameCode:     *game.GameCode,
			ChampionName: *game.ChampionName,
			Game:         string(calculatedGameJSON),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
