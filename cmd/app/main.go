package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Puker228/pingulina/internal/checker"
	"github.com/Puker228/pingulina/internal/runner"
	"github.com/Puker228/pingulina/internal/sites"
)

const (
	userAgent = "pingulina/1.0 (+https://github.com/Puker228/pingulina)"
)

func main() {
	workers := flag.Int("w", 8, "the number of workers")
	websites := flag.String("s", "sites.txt", "file with websites")
	doneFile := flag.String("o", "check.log", "file with output")
	flag.Parse()

	if *workers < 1 {
		log.Fatal("number of workers must be >= 1")
	}

	urls, err := sites.LoadURLs(*websites)
	if err != nil {
		log.Fatal(err)
	}

	logF, err := os.OpenFile(*doneFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	c := checker.NewChecker(20*time.Second, *workers, userAgent)
	r := runner.NewRunner(c, *workers, logF)

	r.Run(urls)

	if err := logF.Close(); err != nil {
		log.Fatal(err)
	}
}
