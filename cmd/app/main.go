package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
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

	_, err = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, time.Since(start), err
}

func loadURLs(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if u := strings.TrimSpace(scanner.Text()); u != "" {
			urls = append(urls, u)
		}
	}
	return urls, scanner.Err()
}

func main() {
	urls, err := loadURLs("sites.txt")
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan string)
	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			for site := range jobs {
				code, dur, err := check(site)
				if err != nil {
					log.Printf("site: %s, error: %v, dur: %v", site, err, dur)
					continue
				}
				log.Printf("site: %s, code: %d, dur: %v", site, code, dur)
			}
		})
	}

	for _, u := range urls {
		jobs <- u
	}
	close(jobs)
	wg.Wait()
}
