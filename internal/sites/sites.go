package sites

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func LoadURLs(path string) ([]string, error) {
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
