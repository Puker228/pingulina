package checker

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Checker struct {
	client    *http.Client
	userAgent string
	maxBody   int64
}

type Result struct {
	Code int
	Dur  time.Duration
	Err  error
}

func NewChecker(timeout time.Duration, maxConns int, userAgent string) *Checker {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.MaxIdleConnsPerHost = maxConns

	return &Checker{
		client:    &http.Client{Timeout: timeout, Transport: tr},
		userAgent: userAgent,
		maxBody:   64 << 10,
	}
}

func (c *Checker) Check(ctx context.Context, site string) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, site, nil)
	if err != nil {
		return Result{0, 0, err}
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return Result{0, time.Since(start), err}
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, io.LimitReader(resp.Body, c.maxBody))
	return Result{resp.StatusCode, time.Since(start), err}
}
