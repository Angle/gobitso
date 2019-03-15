package bitso

import "github.com/shopspring/decimal"

type BookCode string

const (
	// MXN Markets
	BookCode_BTC_MXN 	BookCode = "btc_mxn"
	BookCode_ETH_MXN 	BookCode = "eth_mxn"
	BookCode_XRP_MXN 	BookCode = "xrp_mxn"
	BookCode_LTC_MXN 	BookCode = "ltc_mxn"
	BookCode_BCH_MXN 	BookCode = "bch_mxn"
	BookCode_TUSD_MXN 	BookCode = "tusd_mxn"
	BookCode_MANA_MXN 	BookCode = "mana_mxn"
	BookCode_GNT_MXN 	BookCode = "gnt_mxn"
	BookCode_BAT_MXN 	BookCode = "bat_mxn"

	// BTC Markets
	BookCode_ETH_BTC 	BookCode = "eth_btc"
	BookCode_XRP_BTC 	BookCode = "xrp_btc"
	BookCode_LTC_BTC 	BookCode = "ltc_btc"
	BookCode_BCH_BTC 	BookCode = "bch_btc"
	BookCode_TUSD_BTC 	BookCode = "tusd_btc"
	BookCode_MANA_BTC 	BookCode = "mana_btc"
	BookCode_GNT_BTC 	BookCode = "gnt_btc"
	BookCode_BAT_BTC 	BookCode = "bat_btc"
)

type Book struct {
	BookCode 	BookCode
	Major 		Currency
	Minor 		Currency

	MinimumAmount 	decimal.Decimal // units: major
	MaximumAmount 	decimal.Decimal
	MinimumPrice 	decimal.Decimal // units: minor
	MaximumPrice 	decimal.Decimal
	MinimumValue 	decimal.Decimal // units: minor
	MaximumValue 	decimal.Decimal
}

func FindBookFromTwoCurrencies(books map[BookCode]Book, a, b CurrencyCode) (Book, bool) {
	for _, book := range books {
		if book.Major.Code == a && book.Minor.Code == b {
			return book, true
		} else if book.Minor.Code == a && book.Major.Code == b {
			return book, true
		}
	}

	return Book{}, false
}