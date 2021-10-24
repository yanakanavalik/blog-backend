package main

import (
	"context"
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

	fmt.Fprintln(w, "Articles:")
	for _, v := range visits {
		fmt.Fprintf(w, "[%s] \n", v.Title)
		fmt.Fprintf(w, "%v", v.Paragraph)
	}
	fmt.Fprintln(w, "\nSuccessfully stored an entry of the current request.")
}

type visit struct {
	Timestamp time.Time
	Title     string
	Paragraph []string
	UserIP    string
}

func recordVisit(ctx context.Context, now time.Time) error {
	v := &visit{
		Title:     "First Article",
		Timestamp: time.Now(),
		Paragraph: []string{"First Paragrahla lalal alala", "Second Paragrahla lalal alala", "Third Paragrahla lalal alala", "Fourth Paragrahla lalal alala"},
	}

	k := datastore.IncompleteKey("Article", nil)

	_, err := datastoreClient.Put(ctx, k, v)
	return err
}

func queryVisits(ctx context.Context, limit int64) ([]*visit, error) {
	// Print out previous visits.
	q := datastore.NewQuery("Article").
		Order("-Timestamp").
		Limit(100)

	visits := make([]*visit, 0)
	_, err := datastoreClient.GetAll(ctx, q, &visits)
	return visits, err
}
