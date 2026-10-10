package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const workers = 8

var client = &http.Client{
	Timeout: 2 * time.Minute,
	Transport: &http.Transport{
		MaxIdleConns:        workers,
		MaxIdleConnsPerHost: workers,
		IdleConnTimeout:     90 * time.Second,
	},
}

func check(url string) (int, time.Duration, error) {
	start := time.Now()

	resp, err := client.Get(url)
	if err != nil {
		return 0, time.Since(start), err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	return resp.StatusCode, time.Since(start), err
}

func loadURLs() ([]string, error) {
	file, err := os.Open("sites.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}
		urls = append(urls, url)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func main() {
	data, err := loadURLs()
	if err != nil {
		log.Print(err)
	}

	for _, site := range data {
		code, dur, err := check(site)
		if err != nil {
			log.Print(err)
		}
		log.Printf("site: %v, code: %v, dur: %v", site, code, dur)
	}
}
