package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	workers   = 8
	doneFile  = "check.log"
	userAgent = "pingulina/1.0 (+https://github.com/Puker228/pingulina)"
)

var client = newClient()

func newClient() *http.Client {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.MaxIdleConnsPerHost = workers

	return &http.Client{
		Timeout:   20 * time.Second,
		Transport: tr,
	}
}

func check(site string) (int, time.Duration, error) {
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, site, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return 0, time.Since(start), err
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
	return resp.StatusCode, time.Since(start), err
}

func loadURLs(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for line := 1; scanner.Scan(); line++ {
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" || strings.HasPrefix(raw, "#") {
			continue
		}

		u, err := normalizeURL(raw)
		if err != nil {
			log.Printf("%s:%d: skip %q: %v", path, line, raw, err)
			continue
		}
		if seen[u] {
			continue
		}
		seen[u] = true
		urls = append(urls, u)
	}
	return urls, scanner.Err()
}

func normalizeURL(raw string) (string, error) {
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("unsupported scheme %q", u.Scheme)
	}
	if u.Host == "" {
		return "", errors.New("empty host")
	}
	return u.String(), nil
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
