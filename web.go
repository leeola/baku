package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A basic web server access to a Searcher
//
// TODO: This whole WebListen section needs to be rewritten into a sane
// implementation. The WebListen function and the embedded http handler
// was just a rushed implementation for a working server.
func WebListen(port int, s Searcher, l Logger) {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		var items []SearchItem
		var err error

		c := r.URL.Query().Get("callback")
		q := r.URL.Query().Get("query")
		if q != "" {
			items, err = s.Search(q)
			if err != nil {
				l.Error("Error Searching:", err)
			}
		}

		b, err := json.Marshal(items)
		if err != nil {
			l.Error("Error Marshalling []SearchItem:", err)
		}

		w.Header().Set("Content-Type", "application/json")
		if c != "" {
			fmt.Fprintf(w, "%s(%s)", c, b)
		} else {
			w.Write(b)
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
