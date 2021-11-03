package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

var datastoreClient *datastore.Client

func startDatastore() {
	projectID := os.Getenv("GCLOUD_DATASET_ID")

	ctx := context.Background()

	var err error
	datastoreClient, err = datastore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatal(err)
	}
}

func handleArticlesRequest(w http.ResponseWriter) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Get a list of the most recent visits.
	visits, err := queryVisits(ctx, 10)
	if err != nil {
		msg := fmt.Sprintf("Could not get recent visits: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Record this visit.
	if err := recordVisit(ctx, time.Now()); err != nil {
		msg := fmt.Sprintf("Could not save visit: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	v := &ResponseVisit{
		VisitsArray: visits,
		Count:       len(visits),
	}

	body, err := json.Marshal(v)
	if err != nil {
		msg := fmt.Sprintf("Could not get recent visits: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "[%s]", body)
}

type ResponseVisit struct {
	VisitsArray []*Visit
	Count       int
}

type Visit struct {
	DateCreated time.Time
	Title       string
	Paragraphs  []*Paragraph
}

type Paragraph struct {
	Title string
	Text  []string
}

func recordVisit(ctx context.Context, now time.Time) error {
	v := &Visit{
		Title:       "Welcome!",
		DateCreated: time.Now(),
		Paragraphs: []*Paragraph{
			{
				Title: "",
				Text:  []string{"Welcome to my blog. Nice to meet you", "bla"},
			},
			{
				Title: "Subtitle 1",
				Text:  []string{"Lallalalala 1"},
			},
			{
				Title: "Subtitle 2",
				Text:  []string{"A:A:A::A:"},
			},
		},
	}

	k := datastore.IncompleteKey("Article", nil)

	_, err := datastoreClient.Put(ctx, k, v)
	return err
}

func queryVisits(ctx context.Context, limit int64) ([]*Visit, error) {
	// Print out previous visits.
	q := datastore.NewQuery("Article").
		Order("-DateCreated").
		Limit(100)

	visits := make([]*Visit, 0)
	_, err := datastoreClient.GetAll(ctx, q, &visits)
	return visits, err
}
