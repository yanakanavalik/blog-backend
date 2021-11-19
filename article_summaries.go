package main

import "time"

type ResponseArticleSummary struct {
	ArticlesSummaries []ArticleSummary
	Offset            int
	TotalCount        int
}

type ArticleSummary struct {
	DateCreated time.Time
	Title       string
	Summary     string
	UrlName     string
	Id          int
}
