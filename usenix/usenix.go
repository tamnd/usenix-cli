// Package usenix is the library behind the usenix command line:
// the HTTP client, request shaping, and the typed data models for USENIX
// conference papers scraped from usenix.org.
package usenix

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

const DefaultUserAgent = "Mozilla/5.0 (compatible; usenix-cli/0.1; +https://github.com/tamnd/usenix-cli)"

// Config holds constructor parameters.
type Config struct {
	BaseURL   string
	UserAgent string
	Rate      time.Duration
	Retries   int
	Timeout   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://www.usenix.org",
		UserAgent: DefaultUserAgent,
		Rate:      500 * time.Millisecond,
		Retries:   3,
		Timeout:   30 * time.Second,
	}
}

// Client talks to usenix.org over HTTP.
type Client struct {
	cfg        Config
	httpClient *http.Client
	mu         sync.Mutex
	last       time.Time
}

// NewClient returns a Client with the given config.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

var paperRe = regexp.MustCompile(`href="(/conference/[^/]+/presentation/[^"]+)"[^>]*>([^<]+)</a>`)

// List fetches papers from the given conference's technical-sessions page.
func (c *Client) List(ctx context.Context, conf string, limit int) ([]Paper, error) {
	rawURL := fmt.Sprintf("%s/conference/%s/technical-sessions", c.cfg.BaseURL, conf)
	raw, err := c.get(ctx, rawURL)
	if err != nil {
		return nil, err
	}
	return parsePapers(string(raw), conf, limit, c.cfg.BaseURL), nil
}

// Search fetches all papers from the conference and filters by query (client-side).
func (c *Client) Search(ctx context.Context, query, conf string, limit int) ([]Paper, error) {
	all, err := c.List(ctx, conf, 0)
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(query)
	var out []Paper
	rank := 0
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.Title), q) {
			rank++
			p.Rank = rank
			out = append(out, p)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

func parsePapers(html, conf string, limit int, baseURL string) []Paper {
	matches := paperRe.FindAllStringSubmatch(html, -1)
	var out []Paper
	rank := 0
	for _, m := range matches {
		title := strings.TrimSpace(m[2])
		if title == "" {
			continue
		}
		rank++
		if limit > 0 && rank > limit {
			break
		}
		out = append(out, Paper{
			Rank:       rank,
			Conference: conf,
			Title:      title,
			URL:        baseURL + m[1],
		})
	}
	return out
}

func (c *Client) get(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		b, retry, err := c.do(ctx, rawURL)
		if err == nil {
			return b, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get: %w", lastErr)
}

func (c *Client) do(ctx context.Context, rawURL string) ([]byte, bool, error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}
