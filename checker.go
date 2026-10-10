package main

import (
	"io"
	"net/http"
	"time"
)

type Checker struct {
	client    *http.Client
	userAgent string
	maxBody   int64
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

func (c *Checker) Check(site string) (int, time.Duration, error) {
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, site, nil)
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, time.Since(start), err
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, io.LimitReader(resp.Body, c.maxBody))
	return resp.StatusCode, time.Since(start), err
}
