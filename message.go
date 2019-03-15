package bitso

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

///////////////////////////////
// CATALOG
type ActionType string

const (
	ActionType_NULL 		ActionType = ""
	ActionType_SUBSCRIBE 	ActionType = "subscribe"
)


type Channel string

const (
	Channel_TRADES 			Channel = "trades"
	Channel_DIFF_ORDERS 	Channel = "diff-orders"
	Channel_ORDERS 			Channel = "orders"
	Channel_KEEP_ALIVE 		Channel = "ka"
	Channel_DISCONNECTED 	Channel = "disconnected"
)

type Side int64

const (
	Side_BUY 	Side = 0
	Side_SELL 	Side = 1
)

func (s Side) String() string {
	switch s {
	case Side_BUY:
		return "BUY"
	case Side_SELL:
		return "SELL"
	default:
		return ""
	}
}

///////////////////////////////
////  REST API
// ApiResponse is a general struct used for any response from the REST API, both public and private
type ApiResponse struct {
	Success bool 			`json:"success"`
	Error	ApiError 		`json:"error"`
	Payload	*json.RawMessage `json:"payload"`
}


///////////////////////////////
////  PUBLIC REST API
// Public REST API: Available Books
type PublicAvailableBooksPayload struct {
	Book 			string 			`json:"book"`
	MinimumAmount 	decimal.Decimal `json:"minimum_amount"` // units: major
	MaximumAmount 	decimal.Decimal `json:"maximum_amount"`
	MinimumPrice 	decimal.Decimal `json:"minimum_price"` // units: minor
	MaximumPrice 	decimal.Decimal `json:"maximum_price"`
	MinimumValue 	decimal.Decimal `json:"minimum_value"` // units: minor
	MaximumValue 	decimal.Decimal `json:"maximum_value"`
}

///////////////////////////////
////  PRIVATE REST API
// Private REST API: Account Balance
type PrivateAccountBalancePayload struct {
	Balances 	[]Balance		`json:"balances"`
}

type Balance struct {
	Currency 	CurrencyCode 	`json:"currency"`
	Available 	decimal.Decimal `json:"available"`
	Locked 		decimal.Decimal `json:"locked"`
	Total 		decimal.Decimal `json:"total"`
}

// Private REST API: Account Fees
type PrivateAccountFeesPayload struct {
	Fees			[]Fee	`json:"fees"`
	WithdrawalFees 	map[CurrencyCode]decimal.Decimal `json:"withdrawal_fees"`
}

type Fee struct {
	BookCode		BookCode		`json:"book"`
	TakerFeeDecimal decimal.Decimal `json:"taker_fee_decimal"`
	TakerFeePercent decimal.Decimal `json:"taker_fee_percent"`
	MakerFeeDecimal decimal.Decimal `json:"maker_fee_decimal"`
	MakerFeePercent	decimal.Decimal `json:"maker_fee_percent"`
}


///////////////////////////////
////  WEBSOCKET API
// IncomingMessages is a general struct used for any incoming messages in the websocket feed
type IncomingMessage struct {
	Action  	ActionType 			`json:"action"`
	Channel 	Channel 			`json:"type"`
	Book		BookCode			`json:"book"`
	Sequence 	int64				`json:"sequence"`
	Payload 	*json.RawMessage 	`json:"payload"`
}

// FeddMessage is the actual message type to be passed down to the listener (WebsocketListener)
type FeedMessage struct {
	Channel 	Channel
	Book		BookCode
	Sequence 	int64
	Payload 	interface{}
}


// WebSocket API: Orders Channel
type Orders struct {
	Bids []Offer `json:"bids"`
	Asks []Offer `json:"asks"`
}

type Offer struct {
	Rate 		decimal.Decimal `json:"r"` // units: minor. number of minors per 1 major
	Amount 		decimal.Decimal `json:"a"` // units: major
	Value 		decimal.Decimal `json:"v"` // units: minor
	Side 		Side 			`json:"t"`
	UnixMillis 	int64 			`json:"d"`
}

// Websocket API: Trades Channel
type Trade struct {
	Folio 			int64 			`json:"i"`
	Amount 			decimal.Decimal `json:"a"` // units: major
	Rate 			decimal.Decimal `json:"r"` // units: minor
	Value 			decimal.Decimal `json:"v"` // units: minor
	Side 			Side 			`json:"t"` // maker side
	MakerOrderId 	string 			`json:"mo"`
	TakerOrderId 	string 			`json:"to"`
}


// Websocket API: Subscribe Action Messages
type SubscribeRequestMessage struct {
	Action 	ActionType 		`json:"action"`
	Book 	BookCode 		`json:"book"`
	Channel Channel 		`json:"type"`
}

type SubscribeResponseMessage struct {
	Action 		ActionType 	`json:"action"`
	Response 	string 		`json:"response"`
	Time 		int64 		`json:"time"`
	Channel 	Channel 	`json:"type"`
}