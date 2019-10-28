package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var logFileName = "requests.log"

func logRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			now := time.Now().Unix()
			log.Print(now)
			handler.ServeHTTP(w, r)
		},
	)
}

func parseLines(fileName string) []string {
	logFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	var logs []string
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}
	return logs
}

func countRequests(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Unix()
	numberOfRequests := 0
	logs := parseLines(logFileName)
	re := regexp.MustCompile(`\d{10}`)

	for _, l := range logs {
		match := re.FindStringSubmatch(l)
		if len(match) != 0 {
			matchInt, err := strconv.ParseInt(match[0], 10, 64)
			if err != nil {
				panic(err)
			}
			if now-matchInt <= 60 {
				numberOfRequests = numberOfRequests + 1
			}
		}
	}

	fmt.Fprintln(w, numberOfRequests)
}

func main() {
	port := ":8080"
	handler := http.NewServeMux()
	server := http.Server{Addr: port, Handler: logRequests(handler)}
	logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	handler.HandleFunc("/", countRequests)

	fmt.Println("Server listening at port", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
