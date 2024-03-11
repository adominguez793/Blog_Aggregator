package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/adominguez793/Blog_Aggregator/internal/database"
	"github.com/google/uuid"
)

type Client struct {
	httpClient http.Client
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting feeds every %s on %d goroutines...", timeBetweenRequest, concurrency)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch")
			continue
		}
		log.Printf("Found %d feeds to fetch!", len(feeds))

		waitGroup := &sync.WaitGroup{}
		for _, feed := range feeds {
			waitGroup.Add(1)
			go FeedScrape(db, waitGroup, feed)
		}
		waitGroup.Wait()
	}
}

func FeedScrape(db *database.Queries, waitGroup *sync.WaitGroup, feed database.Feed) {
	defer waitGroup.Done()

	feed, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}
	feedData, err := FeedFetch(feed.Url)
	if err != nil {
		fmt.Printf("Failed to fetch feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Items {
		log.Println("Found Post", item.Title)

		pubDateLayout := "Mon, 02 Jan 2006 15:04:05 -0700"
		pubDate, err := time.Parse(pubDateLayout, item.PubDate)
		if err != nil {
			fmt.Printf("Failed to parse publication date (%s): %v", item.PubDate, err)
			return
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		})
		if err != nil {
			duplicateKeyValueError := "pq: duplicate key value violates unique constraint \"posts_url_key\""
			if !strings.Contains(err.Error(), duplicateKeyValueError) {
				log.Println(err)
			}
		}
	}
	log.Printf("Feed %s collected, %v posts found.", feed.Name, len(feedData.Channel.Items))
}

func FeedFetch(URL string) (*RSS, error) {
	client := NewClient()

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return &RSS{}, err
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return &RSS{}, err
	}
	if resp.StatusCode > 299 {
		return &RSS{}, fmt.Errorf("status code too high: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSS{}, err
	}

	var feedRSS RSS
	err = xml.Unmarshal(dat, &feedRSS)
	if err != nil {
		return &RSS{}, err
	}

	return &feedRSS, nil
}
