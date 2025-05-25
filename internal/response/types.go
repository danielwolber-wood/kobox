package response

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func WriteEpubResponse(w http.ResponseWriter, statusCode int, epubData []byte, filename string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/epub+zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(epubData)))
	w.Write(epubData)
}
