package datstore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

func handleCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Accept, Content-Type")
}

func executeIntParams(r *http.Request, paramName string, defaultValue int) int {
	param, err := strconv.Atoi(r.URL.Query().Get(paramName))
	if err != nil {
		param = defaultValue
	}

	return param
}

func handleArticlesSummariesRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	handleCORS(w)

	articlesStartNum := executeIntParams(r, "start", 0)
	limit := executeIntParams(r, "limit", 10)

	articleSummaries, err := queryArticlesSummaries(ctx, limit, articlesStartNum)
	if err != nil {
		msg := fmt.Sprintf("Could not get article summaries: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	totalArticlesSummariesCount, err := getArticlesCount(ctx)

	v := &ResponseArticleSummary{
		ArticlesSummaries: articleSummaries,
		Offset:            articlesStartNum,
		TotalCount:        totalArticlesSummariesCount,
	}

	body, err := json.Marshal(v)
	if err != nil {
		msg := fmt.Sprintf("Could not json articles summaries: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "[%s]", body)
}

func recordArticleToDataStore(w http.ResponseWriter) {
	ctx := context.Background()

	title := "New Article!"
	dateCreated := time.Now()
	urlName := fmt.Sprintf("article-summary-%s", dateCreated.Format("01-02-2006"))

	if err := recordArticle(ctx, dateCreated, urlName, title); err != nil {
		msg := fmt.Sprintf("Could not record article: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if err := recordArticleSummary(ctx, dateCreated, urlName, title); err != nil {
		msg := fmt.Sprintf("Could not record article summary: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func recordArticle(ctx context.Context, now time.Time, urlName string, title string) error {
	v := &Article{
		Title:       title,
		DateCreated: now,
		UrlName:     urlName,
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

func recordArticleSummary(ctx context.Context, now time.Time, urlName string, title string) error {
	v := &ArticleSummary{
		Title:       title,
		DateCreated: now,
		Summary:     "Article summary. Summary. Summary.",
		UrlName:     urlName,
	}

	k := datastore.IncompleteKey("ArticleSummary", nil)

	_, err := datastoreClient.Put(ctx, k, v)
	return err
}

func queryArticles(ctx context.Context, limit int, start int) ([]*Article, error) {
	q := datastore.NewQuery("Article").
		Order("-DateCreated").
		Limit(limit).
		Offset(start)

	articles := make([]*Article, 0)
	_, err := datastoreClient.GetAll(ctx, q, &articles)
	return articles, err
}

func queryArticlesSummaries(ctx context.Context, limit int, start int) ([]*ArticleSummary, error) {
	q := datastore.NewQuery("ArticleSummary").
		Order("-DateCreated").
		Limit(limit).
		Offset(start)

	articlesSummaries := make([]*ArticleSummary, 0)
	_, err := datastoreClient.GetAll(ctx, q, &articlesSummaries)
	return articlesSummaries, err
}

func getArticlesCount(ctx context.Context) (int, error) {
	q := datastore.NewQuery("ArticleSummary")

	articlesSummaries := make([]*ArticleSummary, 0)
	_, err := datastoreClient.GetAll(ctx, q, &articlesSummaries)
	return len(articlesSummaries), err
}
