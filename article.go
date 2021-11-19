package main

import "time"

type Article struct {
	DateCreated time.Time
	Title       string
	UrlName     string
	Paragraphs  []*Paragraph
	Summary     string
	ID          int
}

type Paragraph struct {
	Title string
	Text  []string
}
