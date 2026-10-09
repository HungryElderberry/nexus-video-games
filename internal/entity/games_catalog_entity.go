package entity

type GamesCatalog struct {
	ExternalGameID string `gorm:"type:varchar(100);primaryKey" json:"external_game_id"`
	Title          string `gorm:"type:varchar(255);not null" json:"title"`
	PriceInBalance int64  `gorm:"type:bigint;not null" json:"price_in_balance"`
	IsAvailable    bool   `gorm:"type:boolean;default:true;not null" json:"is_available"`
}

func (GamesCatalog) TableName() string {
	return "games_catalog"
}
