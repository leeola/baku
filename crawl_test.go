package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const testXmlOneItem string = `
<?xml version="1.0" encoding="utf-8" ?>
<rss
  xmlns:content="http://purl.org/rss/1.0/modules/content/"
  xmlns:wfw="http://wellformedweb.org/CommentAPI/"
  xmlns:dc="http://purl.org/dc/elements/1.1/"
  xmlns:atom="http://www.w3.org/2005/Atom"
  version="2.0">
<channel>
  <title>channel title</title>
  <atom:link
    href="example.com/rss.xml"
    rel="self" type="application/rss+xml" />
  <link>example.com</link>
  <description>channel description</description>
  <language>en</language>
  <item>
    <title>the title</title>
    <link>example.com</link>
    <pubDate>"1970-01-01T00:00:00.000Z"</pubDate>
    <guid isPermaLink="true">the guid</guid>
    <author>the author</author>
    <description>the description</description>
  </item>
</channel>
</rss>
`

const testXmlTwoItems string = `
<?xml version="1.0" encoding="utf-8" ?>
<rss
  xmlns:content="http://purl.org/rss/1.0/modules/content/"
  xmlns:wfw="http://wellformedweb.org/CommentAPI/"
  xmlns:dc="http://purl.org/dc/elements/1.1/"
  xmlns:atom="http://www.w3.org/2005/Atom"
  version="2.0">
<channel>
  <title>channel title</title>
  <atom:link
    href="example.com/rss.xml"
    rel="self" type="application/rss+xml" />
  <link>example.com</link>
  <description>channel description</description>
  <language>en</language>
  <item>
    <title>the title one</title>
    <link>example.com/one</link>
    <pubDate>"1970-01-01T00:00:00.000Z"</pubDate>
    <guid isPermaLink="true">the guid one</guid>
    <author>the author one</author>
    <description>the description one</description>
  </item>
  <item>
    <title>the title two</title>
    <link>example.com/two</link>
    <pubDate>"1970-01-01T00:00:00.000Z"</pubDate>
    <guid isPermaLink="true">the guid two</guid>
    <author>the author two</author>
    <description>the description two</description>
  </item>
</channel>
</rss>
`

func TestFeedCrawlercrawl(t *testing.T) {
	Convey("Should return a lonely feed item", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, testXmlOneItem)
			}))
		defer ts.Close()

		c, _ := NewFeedCrawler(ts.URL)
		fi, err := c.crawl()
		So(err, ShouldBeNil)
		So(len(fi), ShouldEqual, 1)
		So(fi[0], ShouldResemble, XMLFeedItem{
			Title:       "the title",
			Link:        "example.com",
			PubDate:     `"1970-01-01T00:00:00.000Z"`,
			Guid:        "the guid",
			Author:      "the author",
			Description: "the description",
		})
	})

	Convey("Should return all feed items", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, testXmlTwoItems)
			}))
		defer ts.Close()

		c, _ := NewFeedCrawler(ts.URL)
		fi, err := c.crawl()
		So(err, ShouldBeNil)
		So(len(fi), ShouldEqual, 2)
	})
}
