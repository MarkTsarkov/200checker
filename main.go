package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	MakeRequest()
}

func MakeRequest() {
	fmt.Println("Web-site to check status:")

	reader := bufio.NewReader(os.Stdin)
	url, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	url, _ = strings.CutSuffix(url, "\n")
	CheckHttpsPrefix(&url)

	// trunk-ignore(gokart/CWE-918:-Server-Side-Request-Forgery)
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

func CheckHttpsPrefix(url *string) {
	if !strings.HasPrefix(*url, "https://") || !strings.HasPrefix(*url, "http://") {
		*url = "http://" + *url
	}
}
