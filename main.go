package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	//	"net/url"
	"os"
	"strings"
)

func main() {
	var wg sync.WaitGroup

	fmt.Println("Web-sites to check status:")

	scanner := bufio.NewScanner(os.Stdin)
	
	var lines []string
	for {
		scanner.Scan()
		line := scanner.Text()

		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}
	wg.Add(len(lines))

	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	checker := func (url string) {
		defer wg.Done() //уведомление о завершении горутины
		MakeRequest(&url)
	}

	for _, line := range lines {
		go checker(line)
	}

	wg.Wait()
}

func MakeRequest(url *string) {
	curUrl := *url
	CheckHttpsPrefix(&curUrl)
	resp, err := http.Get(curUrl)
	if err != nil {
		log.Fatalln(err)
	}

	r := strings.NewReader(resp.Status)
	body, err := io.ReadAll(r)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(*url + " " + string(body))
}

func CheckHttpsPrefix(url *string) {
	if !strings.HasPrefix(*url, "https://") || !strings.HasPrefix(*url, "http://") {
		*url = "http://www." + *url
	}
}

// trunk-ignore(git-diff-check/error)
// - вынести makeRequest в небольшую функицю
// - в main запускать makeRequest в цикле считывания множества строк с url