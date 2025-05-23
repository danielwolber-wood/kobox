package main

type ReadabilityURLRequest struct {
	Url string `json:"url"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
