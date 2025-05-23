package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/health", handleHealthCheck)
	fmt.Println("Serving on :8080")
	http.ListenAndServe(":8080", r)

}
