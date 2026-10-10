package main

import (
	"fmt"
	"io"
	"log"
	"sync"
)

type Runner struct {
	checker *Checker
	workers int
	out     io.Writer
}

func NewRunner(checker *Checker, workers int, out io.Writer) *Runner {
	return &Runner{
		checker: checker,
		workers: workers,
		out:     out,
	}
}

func (r *Runner) Run(urls []string) {
	jobs := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for range r.workers {
		wg.Go(func() {
			for site := range jobs {
				code, dur, err := r.checker.Check(site)

				var line string
				if err != nil {
					log.Printf("site: %s, error: %v, dur: %v", site, err, dur)
					line = fmt.Sprintf("site: %s, error: %v, dur: %v\n", site, err, dur)
				} else {
					line = fmt.Sprintf("site: %s, code: %d, dur: %v\n", site, code, dur)
				}

				mu.Lock()
				_, err = io.WriteString(r.out, line)
				mu.Unlock()
				if err != nil {
					log.Printf("write result: %v", err)
				}
			}
		})
	}

	for _, u := range urls {
		jobs <- u
	}
	close(jobs)
	wg.Wait()
}
