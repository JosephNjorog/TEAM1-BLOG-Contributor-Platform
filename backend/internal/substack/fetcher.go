package substack

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Fetcher retrieves the publication's published post history. There's no
// general-purpose public API for arbitrary Substack publications, but every
// publication exposes a standard RSS feed at /feed - that's what the live
// implementation reads.
type Fetcher interface {
	FetchPosts(ctx context.Context) ([]Post, error)
}

func NewFetcher(publicationURL string, mock bool) Fetcher {
	if mock {
		return &mockFetcher{}
	}
	return &rssFetcher{publicationURL: strings.TrimSuffix(publicationURL, "/")}
}

type mockFetcher struct{}

func (m *mockFetcher) FetchPosts(_ context.Context) ([]Post, error) {
	now := time.Now()
	return []Post{
		{
			SubstackPostID: "mock-1",
			Title:          "Why Avalanche Subnets Are the Future of App-Specific Chains",
			URL:            "https://team1blog.substack.com/p/mock-subnets-future",
			Author:         "Chidi Contributor",
			PublishedAt:    now.AddDate(0, 0, -45),
		},
		{
			SubstackPostID: "mock-2",
			Title:          "A Field Guide to Avalanche Consensus",
			URL:            "https://team1blog.substack.com/p/mock-consensus-guide",
			Author:         "Chidi Contributor",
			PublishedAt:    now.AddDate(0, 0, -20),
		},
	}, nil
}

type rssFeed struct {
	Channel struct {
		Items []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	Title      string `xml:"title"`
	Link       string `xml:"link"`
	Creator    string `xml:"creator"`
	GUID       string `xml:"guid"`
	PubDateRaw string `xml:"pubDate"`
}

type rssFetcher struct {
	publicationURL string
}

func (f *rssFetcher) FetchPosts(ctx context.Context) ([]Post, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.publicationURL+"/feed", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch substack feed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("substack feed returned status %d", resp.StatusCode)
	}

	var feed rssFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("parse substack feed: %w", err)
	}

	posts := make([]Post, 0, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		publishedAt, _ := time.Parse(time.RFC1123Z, item.PubDateRaw)
		id := item.GUID
		if id == "" {
			id = item.Link
		}
		posts = append(posts, Post{
			SubstackPostID: id,
			Title:          item.Title,
			URL:            item.Link,
			Author:         item.Creator,
			PublishedAt:    publishedAt,
		})
	}
	return posts, nil
}
