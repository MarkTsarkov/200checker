package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	//	"net/url"
	"os"
	"strings"
)

func main() {


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
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

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
	fmt.Println("Processing...")
	for {
		for url := range ch {
			fmt.Println(url)
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
		}
	}
}
	


func CheckHttpsPrefix(url *string) {
	if !strings.HasPrefix(*url, "https://") || !strings.HasPrefix(*url, "http://") {
		*url = "http://"+ string(*url)
	}
}

func UrlInsert(lines []string, ch chan string) {
	fmt.Println("Заливаю урлы в канал...")
	for {
		for i, line := range lines {
			fmt.Printf("Цикл %d прогона началася\n", i)
			ch <- line
			fmt.Printf("В канал залит урл: %s\n", line)
//			if i==len(lines)-1 {
//				i=-1
//			}
			fmt.Printf("Цикл %d прогона закончился\n", i)
		}
		fmt.Println("Урлы залиты, жду обработки...")
		time.Sleep(10*time.Second)
	}


	//ticker := time.NewTicker(5 * time.Second)
	//done := make(chan bool)
//
	//for {
	//	select{
	//	case <-done:
	//		return
	//	case t := <-ticker.C:
	//		fmt.Println("Tick at", t)
	//		fmt.Println("Таймер тикает...")
	//		for _, url := range lines {
	//			ch <- url
	//			fmt.Println("+1 урл в канале...")
	//		}
	//	}
	//	}
	//	
	//time.Sleep(15*time.Second)
	//done <- true
//	_, ok := <- ticker.C; if ok{
//		fmt.Println("Таймер тикает...")
//		for _, url := range lines {
//			ch <- url
//			fmt.Println("+1 урл в канале...")
//		}
//	}
}