package model

type DashboardStatsResponse struct {
	CurrentBalance     int64                     `json:"current_balance"`
	Currency           string                    `json:"currency"`
	TotalGamesOwned    int64                     `json:"total_games_owned"`
	TotalCashSpent     int64                     `json:"total_cash_spent"`
	MonthlyCashSpent   int64                     `json:"monthly_cash_spent"`
	MonthlyGamesPlayed []MonthlyGamePlayResponse `json:"monthly_games_played"`
}

type MonthlyGamePlayResponse struct {
	ExternalGameID  string `json:"external_game_id"`
	Title           string `json:"title"`
	PlaytimeMinutes int    `json:"playtime_minutes"`
}
