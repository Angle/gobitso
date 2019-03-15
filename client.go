package bitso

import (
	"net/http"
	"time"
)

type Client struct {
	// Private API
	key 	string
	secret 	string

	// HTTP Client
	httpClient *http.Client
}

func NewClient() *Client {
	// Initialize the HTTP Client to be used throughout the session
	tr := &http.Transport{
		IdleConnTimeout:    10 * time.Second,
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout: 10 * time.Second,
	}

	return &Client{
		httpClient: httpClient,
	}
}

func (client *Client) SetPrivateKey(key, secret string) {
	client.key = key
	client.secret = secret
}

