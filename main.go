package main

import (
	_ "bufio"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Website struct {
	Company string
	Url 	string
}

func main() {
//	connStr := "user=postgres password=password dbname=checker sslmode=disable"
//	db, err := sql.Open("postgres", connStr)
//	if err != nil {
//        panic(err)
//    } 
//    defer db.Close()

	var lines []Website
	weblist, err := os.Open("web_list.json")
	if err != nil {
		log.Fatal(err)
	}

	defer weblist.Close()

	data, err := io.ReadAll(weblist)
	if err != nil {
		log.Fatal(err)
	}

	jsonErr := json.Unmarshal(data, &lines)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

//	fmt.Println("Web-sites to check status:")
//	scanner := bufio.NewScanner(os.Stdin)
//	for {
//		scanner.Scan()
//		line := scanner.Text()
//
//		if len(line) == 0 {
//			break
//		}
//		lines = append(lines, line)
//	}
//	err := scanner.Err()
//	if err != nil {
//		log.Fatal(err)
//	}


	ch := make(chan string, len(lines))
	go UrlInsert(lines, ch)
	time.Sleep(2*time.Second)
	go MakeRequest(ch)

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig := <-signalChannel
		switch sig {
			case syscall.SIGINT:
			fmt.Println("sigint")
			close(ch)
			return
			case syscall.SIGTERM:
			fmt.Println("sigterm")
			close(ch)
			return
		}
	}
}

func MakeRequest(ch chan string) {
		for url := range ch {
			CheckHttpsPrefix(&url)
			resp, err := http.Get(url)
			if err != nil {
				log.Fatalln(err)
			}


			r := strings.NewReader(resp.Status)
			body, err := io.ReadAll(r)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(url + " " + string(body))

//			dataSlice := strings.Split(string(body), " ")
//			data := make([]string, 3)
//			data[0]=dataSlice[0]+dataSlice[1]
//			data[1]=dataSlice[2]
//			data[2]=dataSlice[3]+dataSlice[4]
//
//			result, err := db.Exec("INSERT INTO logs (Website, Time, Status) values ($1, $2, $3)", 
//			data[0], data[1], data[2])
//			if err != nil{
//				panic(err)
//			}
//			fmt.Println(result.RowsAffected()) 
		}
	
}
	


func CheckHttpsPrefix(url *string) {
	if !strings.HasPrefix(*url, "https://") || !strings.HasPrefix(*url, "http://") {
		*url = "http://"+ string(*url)
	}
}

func UrlInsert(lines []Website, ch chan string) {
	for {
		for _, line := range lines {
			ch <- line.Url
		}
		time.Sleep(10*time.Second)
	}
}

