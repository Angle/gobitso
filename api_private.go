package bitso

import "encoding/json"

func (client *Client) AccountBalance() (map[CurrencyCode]Balance, error) {
	endpoint := "/v3/balance/"

	payload, err := client.httpGet(true, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	rawBalances := PrivateAccountBalancePayload{}
	err = json.Unmarshal(payload, &rawBalances)
	if err != nil {
		return nil, NewHTTPError("cannot parse response payload JSON")
	}

	// Initialize the output map
	m := make(map[CurrencyCode]Balance)

	for _, b := range rawBalances.Balances {
		m[b.Currency] = b
	}

	return m, nil
}


func (client *Client) AccountFees() (map[BookCode]Fee, error) {
	endpoint := "/v3/fees/"

	payload, err := client.httpGet(true, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	rawFees := PrivateAccountFeesPayload{}
	err = json.Unmarshal(payload, &rawFees)
	if err != nil {
		return nil, NewHTTPError("cannot parse response payload JSON")
	}

	// Initialize the output map
	m := make(map[BookCode]Fee)

	for _, f := range rawFees.Fees {
		m[f.BookCode] = f
	}

	return m, nil
}