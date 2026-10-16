package runner

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/Puker228/pingulina/internal/checker"
)

type Runner struct {
	checker *checker.Checker
	workers int
	out     io.Writer
}

func NewRunner(checker *checker.Checker, workers int, out io.Writer) *Runner {
	return &Runner{
		checker: checker,
		workers: workers,
		out:     out,
	}
}

func (r *Runner) Run(ctx context.Context, urls []string) {
	jobs := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for range r.workers {
		wg.Go(func() {
			for site := range jobs {
				res := r.checker.Check(ctx, site)

				var line string
				if res.Err != nil {
					log.Printf("site: %s, error: %v, dur: %v", site, res.Err, res.Dur)
					line = fmt.Sprintf("site: %s, error: %v, dur: %v\n", site, res.Err, res.Dur)
				} else {
					line = fmt.Sprintf("site: %s, code: %d, dur: %v\n", site, res.Code, res.Dur)
				}

				mu.Lock()
				_, res.Err = io.WriteString(r.out, line)
				mu.Unlock()
				if res.Err != nil {
					log.Printf("write result: %v", res.Err)
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
