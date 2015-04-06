package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	CONFIG_FILE = "config.json"
	ES_TAGLINE  = "You Know, for Search"
)

// The config loaded from config.json
type Config struct {
	// The number of *minutes* between crawls of the ToBeIndexed
	CrawlDelay time.Duration `json:"crawl_delay"`
	// The path to use as the bleve index.
	BleveIndex string `json:"bleve_index"`
	// A series of Crawler feeds.
	ToBeIndexed []string `json:"to_be_indexed"`
	// ??
	Elasticsearch string `json:"elasticsearch"`
	// The web server port
	WebPort int `json:"webport"`
}

// A logger interface so that we can change whatever logger we choose.
type Logger interface {
	Debug(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Info(...interface{})
	Print(...interface{})
	Warn(...interface{})
	Debugf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Infof(string, ...interface{})
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
}

func main() {
	// Not sure what logger we want to use. Using this for now because..
	// reasons?
	var l Logger = logrus.New()

	var config, err = parseConfigFile()
	if err != nil {
		l.Errorf("Unable to parse %s, error: %v", CONFIG_FILE, err)
		os.Exit(1)
	}

	os.Exit(runBaku(config, l))
}

func runBaku(config Config, l Logger) (exit int) {
	// Not currently using ElasticSearch. To be implemented very soon.
	//var esUrl = config.Elasticsearch
	//err = checkElasticsearchIsUp(config.Elasticsearch)
	//if err != nil {
	//	log.Printf("Elasticsearch is not reachable at %s, error: %v", esUrl, err)
	//	os.Exit(1)
	//}
	//var failedUrls = checkIndexUrlsAreCrawable(config.ToBeIndexed)
	//if len(failedUrls) > 0 {
	//	log.Printf("%v url(s) are not crawable", failedUrls)
	//	os.Exit(1)
	//}

	s, err := NewBleveSearcher(config.BleveIndex, l)
	if err != nil {
		l.Error(err)
		return 1
	}

	c, err := NewCreepy(config, s, l)
	if err != nil {
		l.Error(err)
		return 1
	}

	// Start the Creepy Crawlers (tehehe)
	c.Start()

	// And the web server
	l.Printf("Web server running on port %d, Ctrl-C to exit.",
		config.WebPort)
	WebListen(config.WebPort, s, l)

	return 0
}

func checkIndexUrlsAreCrawable(urls []string) []string {
	var failedUrls = []string{}
	for _, url := range urls {
		// TODO: Should use a HEAD request here, since the feed's can be quite
		// large. No reason to download them all, i think.
		var resp, err = http.Get(url)
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		resp.Body.Close()

		// Uncommented, because undefined
		//var q Query
		//xml.Unmarshal(xmlFile, &q)

		// log.Println(string(body))
	}

	return failedUrls
}

func parseConfigFile() (Config, error) {
	var config Config

	var file, err = ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func checkElasticsearchIsUp(url string) error {
	var resp, err = http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var esResponse = struct {
		Tagline string `json:"tagline"`
	}{}

	err = json.Unmarshal(body, &esResponse)
	if err != nil {
		return err
	}

	if esResponse.Tagline != ES_TAGLINE {
		var msg = fmt.Sprintf("%v doesn't have tagline we're looking for: %s", url, ES_TAGLINE)
		return errors.New(msg)
	}

	return nil
}
