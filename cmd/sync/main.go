package main

import (
	"fmt"

	"github.com/turpatur/dcg-api/internal/database"
	"github.com/turpatur/dcg-api/internal/services"
)

func main() {
	db := database.ConnectToDB()

	fmt.Println("Starting background card data synchronization...")
	
	services.SyncCards(db)

	fmt.Println("Enforcing card banlist limits...")
	services.ApplyBanlist(db)
}
