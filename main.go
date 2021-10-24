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
	if r.URL.Path == "/Articles" {
		handleArticlesRequest(w)
		return
	}

	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Welcome. Please, use existing routes: \n/Articles")
		return
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}
