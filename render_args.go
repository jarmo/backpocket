package main

import (
	"html/template"
	"net/url"
)

type RenderArgs struct {
	Address     *url.URL
	Title       string
	Image       string
	Excerpt     string
	Byline      string
	SiteName    string
	ReadingTime int
	Content     template.HTML
	ArchivedAt  string
	Error       error
}
