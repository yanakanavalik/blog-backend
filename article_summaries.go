package main

import "time"

type ResponseArticleSummary struct {
	ArticlesSummaries []*ArticleSummary
	Count             int
}

type ArticleSummary struct {
	DateCreated time.Time
	Title       string
	Summary     string
	UrlName     string
}
