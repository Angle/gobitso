package bitso

type CurrencyCode string

const (
	CurrencyCode_BTC 	CurrencyCode = "btc"
	CurrencyCode_MXN 	CurrencyCode = "mxn"
	CurrencyCode_ETH 	CurrencyCode = "eth"
	CurrencyCode_XRP 	CurrencyCode = "xrp"
	CurrencyCode_LTC 	CurrencyCode = "ltc"
	CurrencyCode_BCH 	CurrencyCode = "bch"
	CurrencyCode_TUSD 	CurrencyCode = "tusd"
	CurrencyCode_MANA 	CurrencyCode = "mana"
	CurrencyCode_GNT 	CurrencyCode = "gnt"
	CurrencyCode_BAT 	CurrencyCode = "bat"
)

type Currency struct {
	Code 		CurrencyCode
	Name 		string
	Precision 	int
}

// Go doesn't have constant maps, so we'll make a function that always returns the same map.  This is not optimal for memory allocation.
func CurrencyList() map[CurrencyCode]Currency {
	return map[CurrencyCode]Currency{
		CurrencyCode_BTC: {
			Code: CurrencyCode_BTC,
			Name: "Bitcoin",
			Precision: 8,
		},
		CurrencyCode_MXN: {
			Code: CurrencyCode_MXN,
			Name: "Mexican Pesos",
			Precision: 2,
		},
		CurrencyCode_ETH: {
			Code: CurrencyCode_ETH,
			Name: "Ethereum",
			Precision: 8,
		},
		CurrencyCode_XRP: {
			Code: CurrencyCode_XRP,
			Name: "Ripple",
			Precision: 8,
		},
		CurrencyCode_LTC: {
			Code: CurrencyCode_LTC,
			Name: "Litecoin",
			Precision: 8,
		},
		CurrencyCode_BCH: {
			Code: CurrencyCode_BCH,
			Name: "Bitcoin Cash",
			Precision: 8,
		},
		CurrencyCode_TUSD: {
			Code: CurrencyCode_TUSD,
			Name: "TrueUSD",
			Precision: 2,
		},
		CurrencyCode_MANA: {
			Code: CurrencyCode_MANA,
			Name: "Decentraland",
			Precision: 8,
		},
		CurrencyCode_GNT: {
			Code: CurrencyCode_GNT,
			Name: "Golem",
			Precision: 8,
		},
		CurrencyCode_BAT: {
			Code: CurrencyCode_BAT,
			Name: "Basic Attention Token",
			Precision: 8,
		},
	}
}
