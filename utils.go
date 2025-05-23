package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	errResp := ErrorResponse{
		Status:  statusCode,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		log.Printf("Failed to encode error response: %v", err)
		w.Write([]byte(`{"status":500,"message":"Internal server error"}`))
	}
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, r any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Printf("Failed to encode json response: %v", err)
		w.Write([]byte(`{"status":500,"message":"Internal server error"}`))
		return
	}

}
