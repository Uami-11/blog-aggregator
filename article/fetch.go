// Package article...
package article

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var feed RSSFeed

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return &feed, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return &feed, err
	}

	body := resp.Body
	defer resp.Body.Close()

	decoder := xml.NewDecoder(body)

	if err := decoder.Decode(&feed); err != nil {
		return &feed, err
	}
	fmt.Println("DEBUG right after decode:", feed.Channel.Link)
	feed = *escapeHTML(&feed)

	return &feed, nil
}

func escapeHTML(feed *RSSFeed) *RSSFeed {
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Link = html.UnescapeString(feed.Channel.Link)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
		feed.Channel.Item[i].Link = html.UnescapeString(feed.Channel.Item[i].Link)
	}

	return feed
}
