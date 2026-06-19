package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/turpatur/dcg-api/internal/models"
	"gorm.io/gorm"
)

// SyncCards determines whether to run a Full Sync or Delta Sync and executes it
func SyncCards(db *gorm.DB) {
	var count int64
	if err := db.Model(&models.DBCard{}).Count(&count).Error; err != nil {
		fmt.Println("Error: failed to count cards in database:", err)
		return
	}

	if count == 0 {
		fmt.Println("Empty database detected. Starting full card synchronization...")
		executeFullSync(db)
	} else {
		fmt.Printf("Database contains %d cards. Starting delta card synchronization...\n", count)
		runDeltaSync(db)
	}
}

func fetchMasterList() ([]models.MasterCard, error) {
	fmt.Println("Fetching master card catalog list from API...")
	resp, err := http.Get("https://digimoncard.io/api-public/getAllCards?series=Digimon%20Card%20Game")
	if err != nil {
		return nil, fmt.Errorf("failed to contact master API: %w", err)
	}
	defer resp.Body.Close()

	var masterList []models.MasterCard
	if err := json.NewDecoder(resp.Body).Decode(&masterList); err != nil {
		return nil, fmt.Errorf("failed to decode master card list: %w", err)
	}
	return masterList, nil
}

func executeFullSync(db *gorm.DB) {
	masterList, err := fetchMasterList()
	if err != nil {
		fmt.Println("Error: synchronization failed:", err)
		return
	}

	totalCards := len(masterList)
	fmt.Printf("Found %d cards in the master catalog. Running full synchronization...\n", totalCards)
	syncCardsInChunks(db, masterList)
}

func runDeltaSync(db *gorm.DB) {
	masterList, err := fetchMasterList()
	if err != nil {
		fmt.Println("Error: synchronization failed:", err)
		return
	}

	var existingCards []string
	if err := db.Model(&models.DBCard{}).Pluck("card_number", &existingCards).Error; err != nil {
		fmt.Println("Error: failed to fetch existing cards from database:", err)
		return
	}

	existingSet := make(map[string]bool, len(existingCards))
	for _, code := range existingCards {
		existingSet[code] = true
	}

	var missingCards []models.MasterCard
	for _, card := range masterList {
		if !existingSet[card.CardNumber] {
			missingCards = append(missingCards, card)
		}
	}

	if len(missingCards) == 0 {
		fmt.Println("Database is up-to-date. No new cards to sync.")
		return
	}

	fmt.Printf("Found %d new cards out of %d total catalog cards. Synchronizing missing cards...\n", len(missingCards), len(masterList))
	syncCardsInChunks(db, missingCards)
}

func syncCardsInChunks(db *gorm.DB, cards []models.MasterCard) {
	totalCards := len(cards)
	chunkSize := 40
	var currentChunk []string

	for i, card := range cards {
		currentChunk = append(currentChunk, card.CardNumber)

		if len(currentChunk) == chunkSize || i == totalCards-1 {
			joinedCodes := strings.Join(currentChunk, ",")
			fmt.Printf("Fetching batch details (%d/%d)... ", i+1, totalCards)
			fetchBatchDetails(db, joinedCodes)

			currentChunk = []string{}
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("Data ingestion pipeline completed successfully.")
}

func fetchBatchDetails(db *gorm.DB, cardCodes string) {
	url := "https://digimoncard.io/api-public/search?card=" + cardCodes

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: HTTP request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		fmt.Println("Error: rate limit encountered. Stopping execution to protect IP.")
		return
	}

	var detailedCards []models.DBCard
	if err := json.NewDecoder(resp.Body).Decode(&detailedCards); err != nil {
		fmt.Println("Warning: skipping empty or misformed batch.")
		return
	}

	for _, card := range detailedCards {
		card.DeckLimit = 4
		db.Save(&card)
	}
	fmt.Println("Batch saved to database.")
}
