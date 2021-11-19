package main

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

func handleArticleByIdRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	handleCORS(w)

	articleId := executeIntParams(r, "id", 0)
	article, err := queryArticleById(ctx, articleId)

	if err != nil {
		msg := fmt.Sprintf("Could not find article by id: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(article)
	if err != nil {
		msg := fmt.Sprintf("Could not json article: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "[%s]", body)
}

func recordArticleToDataStore(w http.ResponseWriter) {
	ctx := context.Background()

	title := "New Article!"
	dateCreated := time.Now()

	if err := recordArticle(ctx, dateCreated, "custom-summary-url", title); err != nil {
		msg := fmt.Sprintf("Could not record article: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func recordArticle(ctx context.Context, now time.Time, urlName string, title string) error {
	idNumber, err := getArticlesCount(ctx)
	newId := idNumber + 1

	id := datastore.IDKey("Article", int64(newId), nil)

	v := &Article{
		Title:       title,
		DateCreated: now,
		Summary:     "Article summary. Summary. Summary.",
		UrlName:     urlName,
		ID:          newId,
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

	_, err = datastoreClient.Put(ctx, id, v)
	return err
}

func queryArticleById(ctx context.Context, id int) (*Article, error) {
	key := datastore.IDKey("Article", int64(id), nil)

	article := new(Article)
	err := datastoreClient.Get(ctx, key, article)

	return article, err
}

func queryArticlesSummaries(ctx context.Context, limit int, start int) ([]ArticleSummary, error) {
	q := datastore.NewQuery("Article").
		Order("-DateCreated").
		Limit(limit).
		Offset(start)

	articles := make([]*Article, 0)
	_, err := datastoreClient.GetAll(ctx, q, &articles)

	articlesSummaries := make([]ArticleSummary, len(articles))

	fmt.Print(articles)

	for i := 0; i < len(articles); i++ {
		articlesSummaries[i] = ArticleSummary{
			Id:          articles[i].ID,
			DateCreated: articles[i].DateCreated,
			Summary:     articles[i].Summary,
			Title:       articles[i].Title,
			UrlName:     articles[i].UrlName,
		}

		fmt.Print(articlesSummaries[i])

	}

	fmt.Print(articlesSummaries)

	return articlesSummaries, err
}

func getArticlesCount(ctx context.Context) (int, error) {
	q := datastore.NewQuery("Article")

	articles := make([]*Article, 0)
	_, err := datastoreClient.GetAll(ctx, q, &articles)
	return len(articles), err
}
