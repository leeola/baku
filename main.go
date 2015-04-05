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
)

const (
	CONFIG_FILE = "config.json"
	ES_TAGLINE  = "You Know, for Search"
)

type Config struct {
	// The number of *minutes* between crawls of the ToBeIndexed
	Delay time.Duration `json:"delay"`

	// A series of Crawler feeds.
	ToBeIndexed []string `json:"to_be_indexed"`

	Port          int    `json:"port"`
	Elasticsearch string `json:"elasticsearch"`
}

func main() {
	var config, err = parseConfigFile()
	if err != nil {
		log.Printf("Unable to parse %s, error: %v", CONFIG_FILE, err)
		os.Exit(1)
	}

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

	s, err := NewBleveSearcher("/tmp/baku.bleve.index")
	if err != nil {
		log.Fatal(err)
	}

	c, err := NewCreepy(config, s)
	if err != nil {
		log.Fatal(err)
	}

	// Start the Creepy Crawlers (tehehe)
	c.Start()

	// No web server yet, so waiting here.
	fmt.Println("Web server running on port 3000, Ctrl-C to exit.")
	WebListen(s)
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
