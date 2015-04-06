package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Crawler interface {
	Crawl() ([]SearchItem, error)
}

// The Creepy manages crawlers based on the config settings. Automating
// their requests and logging their output.
type Creepy struct {
	// The config that this crawler will base its actions off of.
	Config Config
	// The crawlers that this Creepy manages.
	Crawlers []Crawler
	// Our search instance
	Searcher Searcher

	logger Logger
}

func NewCreepy(c Config, s Searcher, l Logger) (*Creepy, error) {
	if c.CrawlDelay == 0 {
		return nil, errors.New(fmt.Sprint(
			"baku.Creepy: CrawlDelay ", c.CrawlDelay, "is disabled to "+
				"prevent accidental loop. Use -1 if you need no delay"))
	} else if c.CrawlDelay == -1 {
		// Apparently the caller *really* wants to run without a delay..
		l.Warn("Config CrawlDelay set to -1, running Creepy without delay")
		c.CrawlDelay = 0
	}

	if len(c.ToBeIndexed) == 0 {
		return nil, errors.New(
			"baku.Creepy: Config needs atleast one to_be_indexed")
	}

	// When we support more crawl formats we will analyze the url to
	// decide which crawler to use. For now though, make them all
	// FeedCrawlers.
	crawlers := make([]Crawler, len(c.ToBeIndexed))
	for i, tbi := range c.ToBeIndexed {
		crawlers[i] = &FeedCrawler{Feed: tbi}
	}

	return &Creepy{
		Config:   c,
		Searcher: s,
		Crawlers: crawlers,
		logger:   l,
	}, nil
}

// Loop through our callers, calling them one at a time.
func (c *Creepy) callCrawlers() {
	for _, crawler := range c.Crawlers {
		si, err := crawler.Crawl()
		if err != nil {
			c.logger.Error(err)
			continue
		}

		// Index the searchitems
		c.Searcher.Index(si...)
	}
}

// Start the creepy crawlers.
// TODO: Store a channel to cancel the loop. But for now we don't care.
func (c *Creepy) Start() {
	go func() {
		for {
			c.logger.Info("Crawling")
			c.callCrawlers()
			time.Sleep(c.Config.CrawlDelay * time.Minute)
		}
	}()
}

// These are vastly simplified for the current needs. I welcome
// anyone to correct this XML stuff to be more generic and/or correct
// to the rss spec.
type XMLFeed struct {
	Channel XMLFeedChannel `xml:"channel"`
}

type XMLFeedChannel struct {
	Items []XMLFeedItem `xml:"item"`
}

type XMLFeedItem struct {
	Author      string `xml:"author"`
	Description string `xml:"description"`
	Guid        string `xml:"guid"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Title       string `xml:"title"`
}

// A crawler that crawls an RSS spec XML feed.
type FeedCrawler struct {
	Feed string
}

func NewFeedCrawler(feed string) (*FeedCrawler, error) {
	return &FeedCrawler{Feed: feed}, nil
}

// Query the feed url and return feed items.
func (c *FeedCrawler) crawl() ([]XMLFeedItem, error) {
	res, err := http.Get(c.Feed)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var f XMLFeed
	err = xml.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}

	return f.Channel.Items, nil
}

func (c *FeedCrawler) Crawl() ([]SearchItem, error) {
	fis, err := c.crawl()
	if err != nil {
		return nil, err
	}

	si := make([]SearchItem, len(fis))
	for i, fi := range fis {
		si[i] = SearchItem{
			// Note that author is missing from SearchItem currently. So we're
			// ignoring it here.
			//Author: fi.Author,

			// Also ignoring Summary, simply because storing them both seems
			// excessive since our XML input only has Description currently.
			//Summary: fi.Description,

			Content:     fi.Description,
			Link:        fi.Link,
			PublishedOn: fi.PubDate,
			Title:       fi.Title,
		}
	}

	return si, nil
}
