package bitso

import (
	"encoding/json"
	"errors"
	"strings"
)

// https://bitso.com/api_info#available-books
func (client *Client) AvailableBooks() (map[BookCode]Book, error) {
	endpoint := "/v3/available_books/"

	payload, err := client.httpGet(false, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	//fmt.Println(payload)

	// Parse the response body
	rawBooks := make([]PublicAvailableBooksPayload, 0)
	err = json.Unmarshal(payload, &rawBooks)

	//log.Printf("Got %d books!", len(rawBooks))
	//log.Println(rawBooks)

	// Pull the "constant" currency list map
	currencyList := CurrencyList()

	books := make(map[BookCode]Book)

	for _, rawBook := range rawBooks {
		// Check the book code, determine the currencies from it
		parts := strings.Split(rawBook.Book, "_")

		if len(parts) != 2 {
			// The book code pulled from the API is not in a valid format "major_minor"
			return nil, errors.New("invalid book code string, expecting 'maj_min', got: " + rawBook.Book)
		}

		majorStr := parts[0]
		minorStr := parts[1]

		// Check that both currencies are registered in this library
		if _, exists := currencyList[CurrencyCode(majorStr)]; !exists {
			// major currency code is not registered
			return nil, errors.New("invalid major currency code '" + majorStr + "'")
		}
		if _, exists := currencyList[CurrencyCode(minorStr)]; !exists {
			// minor currency code is not registered
			return nil, errors.New("invalid minor currency code '" + minorStr + "'")
		}

		major := currencyList[CurrencyCode(majorStr)]
		minor := currencyList[CurrencyCode(minorStr)]

		// Create a book from the rawBook
		book := Book {
			BookCode: BookCode(rawBook.Book),
			Major: major,
			Minor: minor,

			MinimumAmount: 	rawBook.MinimumAmount,
			MaximumAmount: 	rawBook.MaximumAmount,
			MinimumPrice: 	rawBook.MinimumPrice,
			MaximumPrice: 	rawBook.MaximumPrice,
			MinimumValue: 	rawBook.MinimumValue,
			MaximumValue: 	rawBook.MaximumValue,
		}

		// Add it to our Book map
		books[book.BookCode] = book
	}

	return books, nil
}
