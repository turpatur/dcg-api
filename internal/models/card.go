package models

import (
	"github.com/google/uuid"
)

// DBCard represents the physical table in Postgres with the full API schema
type DBCard struct {
	CardNumber     string   `gorm:"primaryKey" json:"id"`
	Name           string   `gorm:"not null" json:"name"`
	Type           string   `json:"type"`
	Color          string   `json:"color"`
	Color2         *string  `json:"color2"`
	Level          *int     `json:"level"`
	PlayCost       *int     `json:"play_cost"`
	EvolutionCost  *int     `json:"evolution_cost"`
	EvolutionColor string   `json:"evolution_color"`
	EvolutionLevel *int     `json:"evolution_level"`
	DP             *int     `json:"dp"`
	Attribute      string   `json:"attribute"`
	DigiType       string   `json:"digi_type"`
	DigiType2      *string  `json:"digi_type2"`
	Form           string   `json:"form"`
	Stage          string   `json:"stage"`
	Rarity         string   `json:"rarity"`
	MainEffect     string   `json:"main_effect"`
	SourceEffect   string   `json:"source_effect"`
	AltEffect      string   `json:"alt_effect"`
	
	// SetNames holds set occurrences, serialized as JSON by GORM
	SetNames       []string `gorm:"serializer:json" json:"set_name"`
	DeckLimit      int      `gorm:"default:4" json:"-"`
}

// PairBan maps pair ban rules in the database
type PairBan struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CardA string    `gorm:"not null;index" json:"card_a"`
	CardB string    `gorm:"not null;index" json:"card_b"`
	Note  string    `json:"note"`
}

// MasterCard maps the simplified response from /getAllCards
type MasterCard struct {
	CardNumber string `json:"cardnumber"`
}
