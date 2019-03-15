package main

import (
	"fmt"
	"log"

	"github.com/angle/gobitso"
)

func main() {
	bitsoClient := bitso.NewClient()
	bitsoClient.SetPrivateKey("key", "mysecret")


	// Pull books & limits from the PublicAPI
	log.Println("Pulling books..")
	books, err := bitsoClient.AvailableBooks()
	if err != nil {
		log.Fatalf("Error pulling books from Bitso: %v", err)
	}

	log.Printf("Loaded %d books", len(books))


	// Pull account balance from the Private API
	balances, err := bitsoClient.AccountBalance()
	if err != nil {
		log.Fatalf("Error pulling account balance from Bitso: %v", err)
	}

	fmt.Println("ACCOUNT BALANCES")
	fmt.Println(balances)


	// Pull account fees from the Private API
	fees, err := bitsoClient.AccountFees()
	if err != nil {
		log.Fatalf("Error pulling account fees from Bitso: %v", err)
	}

	fmt.Println("ACCOUNT FEES")
	fmt.Println(fees)


	log.Println("Program ended.")
}