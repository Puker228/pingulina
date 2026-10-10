package main

import (
	"flag"
	"log"
	"os"
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

	if *workers < 1 {
		log.Fatal("number of workers shoud be >= 1")
	}

	urls, err := loadURLs(*sites)
	if err != nil {
		log.Fatal(err)
	}

	logF, err := os.OpenFile(*doneFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	c := NewChecker(20*time.Second, *workers, userAgent)
	r := NewRunner(c, *workers, logF)

	r.Run(urls)

	if err := logF.Close(); err != nil {
		log.Fatal(err)
	}
}
