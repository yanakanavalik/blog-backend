package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	startDatastore()

	http.HandleFunc("/", handleUrl)

	appengine.Main()
}

func handleUrl(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/articles-summaries" {
		handleArticlesSummariesRequest(w, r)
		return
	}

	if r.URL.Path == "/article" {
		handleArticleByIdRequest(w, r)
		return
	}

	if r.URL.Path == "/add-article" {
		recordArticleToDataStore(w)
		return
	}

	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Welcome. Please, use existing routes: \n/articles-summaries -> get existing articles summaries \n /add-article -> add new article (default)")
		return
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}
