package usenix_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tamnd/usenix-cli/usenix"
)

func makeHTML(papers []string) string {
	var sb strings.Builder
	for i, title := range papers {
		sb.WriteString(fmt.Sprintf(`<a href="/conference/test24/presentation/slug%d">%s</a>`, i, title))
	}
	return sb.String()
}

func TestList(t *testing.T) {
	titles := []string{"Paper Alpha", "Paper Beta", "Paper Gamma"}
	html := makeHTML(titles)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(html))
	}))
	defer srv.Close()

	cfg := usenix.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0

	c := usenix.NewClient(cfg)
	papers, err := c.List(context.Background(), "test24", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(papers) != 3 {
		t.Fatalf("got %d papers, want 3", len(papers))
	}
	if papers[0].Title != "Paper Alpha" {
		t.Errorf("title[0] = %q, want %q", papers[0].Title, "Paper Alpha")
	}
	if papers[2].Title != "Paper Gamma" {
		t.Errorf("title[2] = %q, want %q", papers[2].Title, "Paper Gamma")
	}
	wantURL := srv.URL + "/conference/test24/presentation/slug0"
	if papers[0].URL != wantURL {
		t.Errorf("url[0] = %q, want %q", papers[0].URL, wantURL)
	}
	if papers[0].Rank != 1 {
		t.Errorf("rank[0] = %d, want 1", papers[0].Rank)
	}
}

func TestListLimit(t *testing.T) {
	titles := []string{"P1", "P2", "P3", "P4", "P5"}
	html := makeHTML(titles)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(html))
	}))
	defer srv.Close()

	cfg := usenix.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0

	c := usenix.NewClient(cfg)
	papers, err := c.List(context.Background(), "test24", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(papers) != 2 {
		t.Fatalf("got %d papers, want 2 (limit applied)", len(papers))
	}
}

func TestSearch(t *testing.T) {
	titles := []string{"security vulnerability analysis", "network protocol design", "memory safety in C"}
	html := makeHTML(titles)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(html))
	}))
	defer srv.Close()

	cfg := usenix.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0

	c := usenix.NewClient(cfg)
	results, err := c.Search(context.Background(), "net", "test24", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if !strings.Contains(results[0].Title, "network") {
		t.Errorf("expected network paper, got %q", results[0].Title)
	}
}

func TestConferences(t *testing.T) {
	confs := usenix.Conferences()
	if len(confs) < 5 {
		t.Fatalf("got %d conferences, want >= 5", len(confs))
	}
	if confs[0].Year != 2024 {
		t.Errorf("first conf year = %d, want 2024", confs[0].Year)
	}
}
