package bitso

import "fmt"

// Bitso Errors
// https://bitso.com/api_info#error-codes
type ApiError struct {
	Message string `json:"message"`
	Code 	string `json:"code"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("Bitso API Error [%s] %s", e.Code, e.Message)
}


type HTTPError struct {
	msg string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %s", e.msg)
}

func NewHTTPError(m string) HTTPError {
	return HTTPError{
		msg: m,
	}
}


type WebSocketError struct {
	msg string
}

func (e WebSocketError) Error() string {
	return fmt.Sprintf("WebSocket error: %s", e.msg)
}

func NewWebSocketError(m string) WebSocketError {
	return WebSocketError{
		msg: m,
	}
}