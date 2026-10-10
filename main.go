package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	userAgent = "pingulina/1.0 (+https://github.com/Puker228/pingulina)"
)

func main() {
	workers := flag.Int("w", 8, "the number of workers")
	sites := flag.String("s", "sites.txt", "file with websites")
	doneFile := flag.String("o", "check.log", "file with output")
	flag.Parse()

	urls, err := loadURLs(*sites)
	if err != nil {
		log.Fatal(err)
	}

	logF, err := os.OpenFile(*doneFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	c := NewChecker(20*time.Second, *workers, userAgent)

	jobs := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for range *workers {
		wg.Go(func() {
			for site := range jobs {
				code, dur, err := c.Check(site)

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
					log.Printf("write %s: %v", *doneFile, err)
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
