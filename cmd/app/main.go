package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	workers  = 8
	doneFile = "check.log"
)

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

	logF, err := os.OpenFile(doneFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for range workers {
		wg.Go(func() {
			for site := range jobs {
				code, dur, err := check(site)

				var line string
				if err != nil {
					log.Printf("site: %s, error: %v, dur: %v", site, err, dur)
					line = fmt.Sprintf("site: %s, error: %v, dur: %v\n", site, err, dur)
				} else {
					line = fmt.Sprintf("site: %s, code: %d, dur: %v\n", site, code, dur)
				}

				mu.Lock()
				_, err = io.WriteString(logF, line)
				mu.Unlock()
				if err != nil {
					log.Printf("write %s: %v", doneFile, err)
				}
			}
		})
	}

	for _, u := range urls {
		jobs <- u
	}
	close(jobs)
	wg.Wait()

	if err := logF.Close(); err != nil {
		log.Fatal(err)
	}
}
