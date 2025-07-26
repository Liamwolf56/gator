package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html" // For html.UnescapeString
	"io"
	"net/http"
	"time" // For http.Client timeout
)

// RSSFeed represents the top-level structure of an RSS feed XML.
type RSSFeed struct {
	XMLName xml.Name `xml:"rss"` // This field maps to the <rss> root element
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"` // A slice to hold multiple <item> elements
	} `xml:"channel"` // This field maps to the <channel> element
}

// RSSItem represents a single item (e.g., an article or podcast episode) in an RSS feed.
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"` // This will be a string, we won't parse it to time.Time here
}

// fetchFeed fetches an RSS feed from the given URL, parses it,
// and returns an RSSFeed struct.
func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create an HTTP client with a timeout to prevent hanging connections.
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new HTTP GET request with the provided context.
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the User-Agent header to identify our application.
	req.Header.Set("User-Agent", "gator")

	// Execute the HTTP request.
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed from %s: %w", feedURL, err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check for a successful HTTP status code.
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("bad status code from %s: %d", feedURL, resp.StatusCode)
	}

	// Read the entire response body.
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rssFeed RSSFeed
	// Unmarshal the XML data into the RSSFeed struct.
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	// Decode HTML entities in relevant fields for the channel.
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	// Decode HTML entities in relevant fields for each item.
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	return &rssFeed, nil
}
