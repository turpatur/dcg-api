package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/turpatur/dcg-api/internal/models"
	"gorm.io/gorm"
)

// ApplyBanlist enforces official card limits and pair bans in the database
func ApplyBanlist(db *gorm.DB) {
	fmt.Println("Applying standard limits...")

	restrictedTo1 := []string{
		"BT23-032", "BT3-092", "BT10-080", "EX5-059", "EX5-061",
		"BT1-090", "BT6-104", "BT13-110", "BT16-011", "EX3-057",
		"EX4-006", "EX1-021", "BT19-040", "EX2-070", "BT4-111",
		"BT17-069", "BT4-104", "P-029", "P-030", "BT11-033",
		"ST9-09", "EX4-030", "P-123", "P-130", "BT15-057",
		"BT9-098", "ST2-13", "BT14-084", "BT14-002", "BT15-102",
		"EX5-015", "EX5-018", "EX5-062", "BT13-012", "BT2-069",
		"BT7-069", "BT3-054", "EX2-039", "P-008", "P-025",
		"BT11-064", "BT7-107", "BT10-009", "BT7-038", "BT7-064",
		"BT2-047", "BT3-103", "BT6-100", "EX1-068", "BT7-072",
	}

	bannedTo0 := []string{
		"BT5-109", "BT2-090", "EX5-065",
	}

	upTo50 := []string{
		"BT6-085", "BT11-061", "BT22-079", "EX2-046", "EX9-048", "EX11-027",
	}

	db.Model(&models.DBCard{}).Where("card_number IN ?", restrictedTo1).Update("deck_limit", 1)
	db.Model(&models.DBCard{}).Where("card_number IN ?", bannedTo0).Update("deck_limit", 0)
	db.Model(&models.DBCard{}).Where("card_number IN ?", upTo50).Update("deck_limit", 50)

	fmt.Println("Card deck limits applied successfully.")

	fmt.Println("Applying pair ban rulings...")
	pairBans := []models.PairBan{
		{ID: uuid.New(), CardA: "EX2-007", CardB: "EX7-064", Note: "Official Pair Ban"},
		{ID: uuid.New(), CardA: "BT20-037", CardB: "BT17-035", Note: "Official Pair Ban"},
		{ID: uuid.New(), CardA: "BT20-037", CardB: "EX8-037", Note: "Official Pair Ban"},
	}

	db.Exec("TRUNCATE TABLE pair_bans;")
	db.Create(&pairBans)

	fmt.Println("Pair ban rulings applied successfully.")
}
