package main

import "net/http"

func crawlUrl(url string) {
	http.Get(url)
}
