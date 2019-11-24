package model

// Response represents the response of an operation.
type Response struct {
	Entity     interface{} `json:"entity"`
	Errors     []string    `json:"errors"`
	StatusCode int
}
