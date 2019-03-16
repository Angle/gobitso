# Bitso

A pure Go implementation of the [Bitso API v3](https://bitso.com/api_info), including the Public, Private and Websocket APIs.

[Bitso](https://bitso.com) is the leading Mexican cryptocurrency exchange.



This is API is being used in production, even thought at the moment it's missing some functions. We'll keep adding functionality
until we cover all of the official API.

## Upcoming Features
- [ ] Place Limit and Market orders.
- [ ] WebSocket Diff-Order channel multiplexer, with support for persistent connections and auto-recovery.

## How to Use
See the examples at `example_websocket/` and `example_private_api/`.

### WebSocket Listener loop
```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/angle/gobitso"
)

func main() {
	bitsoClient := bitso.NewClient()

	// Pull books & limits from the API
	log.Println("Pulling books..")
	books, err := bitsoClient.AvailableBooks()
	if err != nil {
		log.Fatalf("Error pulling books from Bitso: %v", err)
	}

	log.Printf("Loaded %d books", len(books))


	// INITIALIZE WEBSOCKET BITSO CLIENT
	bitsoWs := bitso.NewWebsocketListener()

	feed, err := bitsoWs.Connect()

	if err != nil {
		log.Fatal("Error connecting to Bitso's websocket")
	}

	// subscribe to _every_ book
	for bookCode := range books {
		err = bitsoWs.Subscribe(bookCode, bitso.Channel_ORDERS)
		if err != nil {
			log.Fatalf("Error subscribing to '%s' Orders channel", bookCode)
		}

		err = bitsoWs.Subscribe(bookCode, bitso.Channel_TRADES)
		if err != nil {
			log.Fatalf("Error subscribing to '%s' Trades channel", bookCode)
		}
	}

	// Catch interrupts (Ctrl-c)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Start main websocket listener loop
	MainLoop:
	for {
		select {
		case msg := <-feed:
			// Incoming message from the websocket!
			switch msg.Channel {
			case bitso.Channel_DISCONNECTED:
				log.Println("Bitso websocket was disconnected!")
				break MainLoop

			case bitso.Channel_ORDERS:
				orders, ok := msg.Payload.(bitso.Orders)

				if !ok {
					// invalid payload in orders
					log.Println("invalid orders payload")
					continue
				}

				// Do something with the incoming orders payload for the book
				out := fmt.Sprintf("[%s] ", msg.Book)


				// Example: print the Best Bid / Ask
				if len(orders.Bids) > 0 {
					out += fmt.Sprintf("BID: %s @ $%s", orders.Bids[0].Amount.StringFixed(8), orders.Bids[0].Rate.StringFixed(8))
				} else {
					out += "BID: --"
				}

				out += " | "

				if len(orders.Asks) > 0 {
					out += fmt.Sprintf("ASK: %s @ $%s", orders.Asks[0].Amount.StringFixed(8), orders.Asks[0].Rate.StringFixed(8))
				} else {
					out += "ASK: --"
				}

				log.Println(out)



			case bitso.Channel_TRADES:
				trades, ok := msg.Payload.([]bitso.Trade)

				if !ok {
					// invalid payload in trades
					log.Println("invalid trades payload")
					continue
				}

				// Do something with the incoming trades payload for the book
				// Example: print the Trade
				for _, trade := range trades {
					log.Printf("[%s] TRADE: [%s] %s @ $%s = $%s", msg.Book, trade.Side.String(), trade.Amount.StringFixed(8), trade.Rate.StringFixed(8), trade.Value.StringFixed( 2))
				}

			}
		case <-interrupt:
			log.Println("interrupt!")
			bitsoWs.Disconnect()
			break MainLoop
		}
	}

	log.Println("Program ended.")
}
```

**Sample output**
```
2019/03/15 16:00:35 Pulling books..
2019/03/15 16:00:35 Loaded 17 books
2019/03/15 16:00:35 BitsoWebsocket: connecting to wss://ws.bitso.com
2019/03/15 16:00:36 BitsoWebsocket: connected!
2019/03/15 16:00:36 BitsoWebsocket: ORDERS subscription ok!
2019/03/15 16:00:36 BitsoWebsocket: TRADES subscription ok!
...
2019/03/15 16:00:42 [btc_mxn] BID: 0.16692328 @ $74200.24000000 | ASK: 0.13650000 @ $74797.11000000
2019/03/15 16:00:42 [eth_btc] BID: 1.52950000 @ $0.03452800 | ASK: 2.60000000 @ $0.03494999
2019/03/15 16:00:42 [xrp_mxn] BID: 1290.32853900 @ $5.93000000 | ASK: 126.92727020 @ $5.96000000
2019/03/15 16:00:42 [bch_mxn] BID: 0.41763786 @ $2681.78000000 | ASK: 12.00000000 @ $2748.08000000
2019/03/15 16:00:42 [bch_mxn] BID: 0.41763786 @ $2681.78000000 | ASK: 12.00000000 @ $2748.08000000
2019/03/15 16:00:43 [xrp_mxn] TRADE: [BUY] 34.89439700 @ $5.93000000 = $206.92
2019/03/15 16:00:43 [ltc_mxn] BID: 0.05123935 @ $1106.00000000 | ASK: 60.00000000 @ $1115.55000000
2019/03/15 16:00:43 [bch_mxn] BID: 0.41763786 @ $2681.78000000 | ASK: 1.96140284 @ $2748.05000000
2019/03/15 16:00:43 [tusd_btc] BID: 0.59000000 @ $0.00025830 | ASK: 3.03000000 @ $0.00026006
2019/03/15 16:00:43 [eth_btc] BID: 1.52950000 @ $0.03452800 | ASK: 2.60000000 @ $0.03494999
2019/03/15 16:00:43 [ltc_mxn] BID: 0.05123935 @ $1106.00000000 | ASK: 7.12927867 @ $1115.52000000
...
^C
2019/03/15 16:00:47 interrupt!
2019/03/15 16:00:47 BitsoWebsocket: quit, attempting clean disconnect
2019/03/15 16:00:47 BitsoWebsocket: read: websocket: close 1000 (normal)
2019/03/15 16:00:48 Program ended.
```

## Functionality
### Public REST API
- [x] Available Books
- [ ] Ticker
- [ ] Order Book
- [ ] Trades 

### Private REST API
- [x] Generating API Keys
- [x] Creating and Signing Requests
- [ ] Account Status
- [ ] Document Upload
- [ ] Mobile Phone Number Registration
- [ ] Mobile Phone Number Verification
- [x] Account Balance
- [x] Fees
- [ ] Ledger
- [ ] Withdrawals
- [ ] Fundings
- [ ] User Trades
- [ ] Order Trades
- [ ] Open Orders
- [ ] Lookup Orders
- [ ] Cancel Order
- [ ] Place an Order
- [ ] Funding Destination
- [ ] Crypto Withdrawals
- [ ] SPEI Withdrawal
- [ ] Bank codes
- [ ] Debit Card Withdrawal
- [ ] Phone Number Withdrawal


### WebSocket API
- [x] Trades Channel
- [ ] Diff-Order Channel
- [x] Orders Channel

## Notes
- Bitso sometimes sends scientific notation for very small numbers (ex. `1E-8`); however, the `shopspring/decimal`
  library automatically handles scientific notation, so we _shouldn't_ have to worry about it.
- Bitso constantly adds and removes currencies in their platform, and while this library can dynamically handle most changes,
  some constants are currently hardcoded in two files: `currency.go` and `book.go`. I'll try to keep it up-to-date, but feel
  free to send a PR with any new currencies or books.

## About
Built by [edmundofuentes](https://github.com/edmundofuentes) for [Angle](https://www.angle.mx).