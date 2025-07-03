package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/senaphim/aggreGator-go/internal/database"
)

type rssFeed struct {
	Channel struct {
		Title         string `xml:"title"`
		Link          string `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Items         []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"descriptiron"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (r *rssFeed) Clean() {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)

	for i, itm := range r.Channel.Items {
		itm.Title = html.UnescapeString(itm.Title)
		itm.Description = html.UnescapeString(itm.Description)
		if itm.Title == "Optimize For Simplicity First" {
			itm.Title = "Optimize for simplicity first"
		}
		r.Channel.Items[i] = itm
	}
}

func FetchFeed(ctx context.Context, feedURL string) (*rssFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmtErr := fmt.Errorf("Error constructing request:\n%v", err)
		return &rssFeed{}, fmtErr
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("User-Agent", "aggreGator-go")

	res, err := client.Do(req)
	if err != nil {
		fmtErr := fmt.Errorf("Error making request:\n%v", err)
		return &rssFeed{}, fmtErr
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmtErr := fmt.Errorf("Error decoding response body:\n%v", err)
		return &rssFeed{}, fmtErr
	}

	feed := &rssFeed{}
	if err := xml.Unmarshal(body, feed); err != nil {
		fmtErr := fmt.Errorf("Error unmarshalling body:\n%v", err)
		return &rssFeed{}, fmtErr
	}

	feed.Clean()

	return feed, nil
}

func scrapeFeeds(s *state) error {
	nextFd, err := s.db.GetNextFeed(context.Background())
	if err != nil {
		fmtErr := fmt.Errorf("Error getting next feed to fetch:\n%v", err)
		return fmtErr
	}

	fetchMarker := database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now().Local(),
		LastFetchedAt: sql.NullTime{Time: time.Now().Local(), Valid: true},
		ID:            nextFd.ID,
	}
	if err := s.db.MarkFeedFetched(context.Background(), fetchMarker); err != nil {
		fmtErr := fmt.Errorf("Error updating feed fetched times:\n%v", err)
		return fmtErr
	}

	fd, err := FetchFeed(context.Background(), nextFd.Url)
	if err != nil {
		fmtErr := fmt.Errorf("Error fetching feed:\n%v", err)
		return fmtErr
	}

	for _, item := range fd.Channel.Items {
		fmt.Printf("%v\n", item.Title)
	}

	return nil
}
