package structs

import (
	"time"
)

type Account struct {
	AccountID       string            `gorm:"column:account_id;primaryKey"`
	Email           string            `gorm:"column:email"`
	CreatedAt       time.Time         `gorm:"column:created_at"`
	Profile         *string           `gorm:"column:profile"`
	AccountSummoner []AccountSummoner `gorm:"foreignKey:AccountID"`
}

type AccountSummoner struct {
	AccountID    string    `gorm:"column:account_id;primaryKey"`
	SummonerName string    `gorm:"column:summoner_name;primaryKey"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	Account      Account   `gorm:"foreignKey:AccountID;references:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Game struct {
	GameID       string     `gorm:"column:game_id;primaryKey"`
	SummonerName *string    `gorm:"column:summoner_name"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	GameCode     *string    `gorm:"column:game_code"`
	ChampionName *string    `gorm:"column:champion_name"`
	GameData     []GameData `gorm:"foreignKey:GameID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type GameData struct {
	ID           int64   `gorm:"column:id;primaryKey;autoIncrement"`
	GameID       string  `gorm:"column:game_id"`
	GameTime     float64 `gorm:"column:game_time"`
	GameData     string  `gorm:"column:game_data"`
	SummonerName *string `gorm:"column:summoner_name"`
	ChampionName *string `gorm:"column:champion_name"`
	Games        Game    `gorm:"foreignKey:GameID;references:GameID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Log struct {
	LogID       int       `gorm:"column:log_id;primaryKey;autoIncrement"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	IP          *string   `gorm:"column:ip"`
	Event       *string   `gorm:"column:event"`
	AccessToken *string   `gorm:"column:access_token"`
	LogStatus   *string   `gorm:"column:log_status"`
	AccountID   *string   `gorm:"column:account_id"`
}

type LastByCodeRequest struct {
	Code string `json:"code"`
	Item string `json:"item"`
}

type HTTPErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LastByCodeResponseData struct {
	GameID       string    `json:"game_id"`
	SummonerName string    `json:"summoner_name"`
	CreatedAt    time.Time `json:"created_at"`
	GameCode     string    `json:"game_code"`
	ChampionName string    `json:"champion_name"`
	Game         string    `json:"game"`
}

type LastByCodeResponse struct {
	Success bool                   `json:"success"`
	Data    LastByCodeResponseData `json:"data"`
}
