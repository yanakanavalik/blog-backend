package main

import "time"

type ResponseArticle struct {
	Articles []*Article
	Count    int
}

type Article struct {
	DateCreated time.Time
	Title       string
	UrlName     string
	Paragraphs  []*Paragraph
}

type Paragraph struct {
	Title string
	Text  []string
}
