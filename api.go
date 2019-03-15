package bitso

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const API_ENDPOINT = "https://api.bitso.com"
// Common API functions


// error code -1 means unknown error (not bitso)
// error code 0 means no error
// error code >0 is a Bitso error
func (client *Client) httpGet(private bool, endpoint string, items []string, query map[string]string) ([]byte, error) {
	u, _ := url.Parse(API_ENDPOINT)
	u.Path = endpoint

	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		// error building the request from the given parameters
	}

	// Add additional required headers headers
	request.Header.Add("Content-type", "application/json")


	// Check if the API Call is Private, if so, add an Authorization header
	if private {
		if client.key == "" || client.secret == "" {
			return []byte(""), NewHTTPError("client's private key/secret pair is not set")
		}

		authHeader := client.buildSignature("GET", endpoint, "")

		// Add custom headers
		request.Header.Add("Authorization", authHeader)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return []byte(""), NewHTTPError(fmt.Sprintf("http request error: %v", err))
	}
	defer response.Body.Close()

	/*
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	*/

	// Attempt to parse payload
	responsePayload, err := parseResponse(response)

	if err != nil {
		return []byte(""), err
	}

	return responsePayload, nil
}


// error code -1 means unknown error (not bitso)
// error code 0 means no error
// error code >0 is a Bitso error
func (client *Client) httpPost(private bool, endpoint string, payload map[string]string) ([]byte, error) {
	u, _ := url.Parse(API_ENDPOINT)
	u.Path = endpoint

	// Convert the Payload to a json string
	payloadString, err := json.Marshal(payload)
	if err != nil {
		//
		return []byte(""), NewHTTPError("could not build a JSON string from the given payload")
	}

	request, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(payloadString))
	if err != nil {
		// error building the request from the given parameters
	}

	// Add additional required headers headers
	request.Header.Add("Content-type", "application/json")


	// Check if the API Call is Private, if so, add an Authorization header
	if private {
		if client.key == "" || client.secret == "" {
			return []byte(""), NewHTTPError("client's private key/secret pair is not set")
		}

		authHeader := client.buildSignature("POST", endpoint, string(payloadString))

		// Add custom headers
		request.Header.Add("Authorization", authHeader)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return []byte(""), NewHTTPError(fmt.Sprintf("http request error: %v", err))
	}
	defer response.Body.Close()

	/*
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	*/

	// Attempt to parse payload
	responsePayload, err := parseResponse(response)

	if err != nil {
		return []byte(""), err
	}

	return responsePayload, nil
}

func (client *Client) buildSignature(method, endpoint, payload string) string {

	// Generate a Nonce from the current nano time
	nonce := strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond), 10)

	// fmt.Println("nonce", nonce)

	// Package hmac implements the Keyed-Hash Message Authentication Code (HMAC) as defined in U.S. Federal
	// Information Processing Standards Publication 198. An HMAC is a cryptographic hash that uses a key to sign
	// a message. The receiver verifies the hash by recomputing it using the same key.

	// Compile the message that should be signed
	message := nonce + method + endpoint + payload

	// fmt.Println("message", message)

	// Initialize the HMAC according to the required Hash and using the client's secret
	h := hmac.New(sha256.New, []byte(client.secret))

	// Write Data to it
	h.Write([]byte(message))

	// Get result and encode as hexadecimal string, this is our signature
	signature := hex.EncodeToString(h.Sum(nil))

	// fmt.Println("signature", signature)

	// Build the header string
	authHeader := fmt.Sprintf("Bitso %s:%s:%s", client.key, nonce, signature)

	// fmt.Println("header", authHeader)

	return authHeader
}

// error -1 means unknown error, could not parse the response body
// error 0 means no error
func parseResponse(response *http.Response) (payload []byte, err error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// cannot even read the response body, error with the reader interface
		return []byte(""), NewHTTPError("cannot read response body, invalid reader interface")
	}

	// debug responde body (raw payload)
	// fmt.Println(string(body))

	// Attempt to parse the body as JSON
	msg := ApiResponse{}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		// could not parse to json
		log.Println(err)
		return []byte(""), NewHTTPError("cannot parse JSON in response body")
	}

	if response.StatusCode != http.StatusOK || msg.Success == false ||  msg.Error.Code != "" {
		// there was some error in the request, pass it down
		return []byte(""), msg.Error
	}

	return *msg.Payload, nil
}