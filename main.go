package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	MakeRequest()
}

func MakeRequest() {
	resp, err := http.Get("https://www.amazon.com/")
	if err != nil {
		log.Fatalln(err)
	}

	r := strings.NewReader(resp.Status)
	body, err := io.ReadAll(r)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}