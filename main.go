package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"xml"
)

const (
	CONFIG_FILE = "config.json"
	ES_TAGLINE  = "You Know, for Search"
)

type Config struct {
	ToBeIndexed   []string `json:"to_be_indexed"`
	Port          int      `json:"port"`
	Elasticsearch string   `json:"elasticsearch"`
}

func main() {
	var config, err = parseConfigFile()
	if err != nil {
		log.Printf("Unable to parse %s, error: %v", CONFIG_FILE, err)
		os.Exit(1)
	}

	var esUrl = config.Elasticsearch

	err = checkElasticsearchIsUp(config.Elasticsearch)
	if err != nil {
		log.Printf("Elasticsearch is not reachable at %s, error: %v", esUrl, err)
		os.Exit(1)
	}

	var failedUrls = checkIndexUrlsAreCrawable(config.ToBeIndexed)
	if len(failedUrls) > 0 {
		log.Printf("%v url(s) are not crawable", failedUrls)
		os.Exit(1)
	}
}

func checkIndexUrlsAreCrawable(urls []string) []string {
	var failedUrls = []string{}
	for _, url := range urls {
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

		var q Query
		xml.Unmarshal(xmlFile, &q)

		// log.Println(string(body))
	}

	return failedUrls
}

func parseConfigFile() (*Config, error) {
	var config Config

	var file, err = ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
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
