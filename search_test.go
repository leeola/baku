package main

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type FakeLogger struct{}

func (l *FakeLogger) Debug(a ...interface{})            {}
func (l *FakeLogger) Debugf(s string, a ...interface{}) {}
func (l *FakeLogger) Error(a ...interface{})            {}
func (l *FakeLogger) Errorf(s string, a ...interface{}) {}
func (l *FakeLogger) Fatal(a ...interface{})            {}
func (l *FakeLogger) Fatalf(s string, a ...interface{}) {}
func (l *FakeLogger) Info(a ...interface{})             {}
func (l *FakeLogger) Infof(s string, a ...interface{})  {}
func (l *FakeLogger) Print(a ...interface{})            {}
func (l *FakeLogger) Printf(s string, a ...interface{}) {}
func (l *FakeLogger) Warn(a ...interface{})             {}
func (l *FakeLogger) Warnf(s string, a ...interface{})  {}

func TestBleveSearcherSearch(t *testing.T) {
	tmpDir := "./_test/tmp"
	os.MkdirAll(tmpDir, 0644)

	Convey("A single index", t, func() {
		indexDir := filepath.Join(tmpDir, "index.bleve")
		// Remove the tmp dir
		os.RemoveAll(indexDir)

		s, err := NewBleveSearcher(indexDir, &FakeLogger{})
		So(err, ShouldBeNil)

		err = s.Index([]SearchItem{
			SearchItem{
				Link:    "the link",
				Title:   "the title",
				Content: "the content",
			},
		}...)
		So(err, ShouldBeNil)

		Convey(`Should match "title"`, func() {
			r, err := s.Search("title")
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 1)
		})

		Convey(`Should match "content"`, func() {
			r, err := s.Search("content")
			So(err, ShouldBeNil)
			So(len(r), ShouldEqual, 1)
		})

		Convey("Should return the match", func() {
			r, err := s.Search("title")
			So(err, ShouldBeNil)
			So(r[0], ShouldResemble, SearchItem{
				Link:    "the link",
				Title:   "the title",
				Content: "the content",
				Summary: "the content",
			})
		})
	})
}
