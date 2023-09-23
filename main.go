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
	var url string
	fmt.Println("Web-site to check status:\n")

	reader := bufio.NewReader(os.Stdin) 
	url, err = reader.ReadString("\n")
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Get(url)
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