package crawler

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func Test_ItReturnsAnErrorIfURLIsUnreachable(t *testing.T) {
	fetcher := &HTTPFetcher{}

	_, err := fetcher.Fetch("http://localhost:40123")

	assert.Contains(t, err.Error(), "unable to fetch url:")
}

func Test_ItScrapesStaticResourcesFromHTML(t *testing.T) {
	fetcher := &HTTPFetcher{}
	body := new(bytes.Buffer)
	pageResults := NewPageResults()

	_, err := body.WriteString(`
	<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
	<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
		<title>Title</title>

		<link rel="stylesheet" href="style.css"/>
		<link rel="stylesheet" href="/static/css/style.css"/>
		<link rel="stylesheet" href="http://www.google.com/style.css"/>
		<style></style>

		<script src="script.js"></script>
		<script src="/static/js/script.js"></script>
		<script src="http://www.google.com/script.js"></script>
		<script></script>
	</head>

	<body>
		<img src="/../avatar.jpg">
	</body>
	</html>`)
	assert.Nil(t, err)

	expectedPageResults := NewPageResults()
	expectedPageResults.cssStaticResources = []string{
		"http://localhost:40123/style.css",
		"http://localhost:40123/static/css/style.css",
		"http://www.google.com/style.css",
	}
	expectedPageResults.jsStaticResources = []string{
		"http://localhost:40123/script.js",
		"http://localhost:40123/static/js/script.js",
		"http://www.google.com/script.js",
	}
	expectedPageResults.imgStaticResources = map[string]bool{
		"http://localhost:40123/avatar.jpg": true,
	}

	queryDoc, err := goquery.NewDocumentFromReader(body)
	assert.Nil(t, err)

	queryDoc.Url, err = url.Parse("http://localhost:40123")
	assert.Nil(t, err)

	fetcher.scrapeStaticResources(queryDoc, pageResults)

	assert.Equal(t, expectedPageResults, pageResults)
}

func Test_ItScrapesLinksFromHTML(t *testing.T) {
	fetcher := &HTTPFetcher{}
	body := new(bytes.Buffer)
	pageResults := NewPageResults()

	_, err := body.WriteString(`
	<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
	<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
		<title>Title</title>

		<link rel="stylesheet" href="style.css"/>
		<link rel="stylesheet" href="/static/css/style.css"/>
		<link rel="stylesheet" href="http://www.google.com/style.css"/>
		<style></style>

		<script src="script.js"></script>
		<script src="/static/js/script.js"></script>
		<script src="http://www.google.com/script.js"></script>
		<script></script>
	</head>

	<body>
		<a href="/">Home</a>
		<a href="/blog/">Blog</a>
		<a href="http://localhost:40123/absolute/">Absolute</a>
		<a href="http://www.google.com">Google</a>
		<a href="http://www.facebook.com">Facebook</a>
	</body>
	</html>`)
	assert.Nil(t, err)

	expectedPageResults := NewPageResults()
	expectedPageResults.internalURLs = map[string]bool{
		"http://localhost:40123/":          true,
		"http://localhost:40123/blog/":     true,
		"http://localhost:40123/absolute/": true,
	}
	expectedPageResults.externalURLs = map[string]bool{
		"http://www.google.com":   true,
		"http://www.facebook.com": true,
	}

	queryDoc, err := goquery.NewDocumentFromReader(body)
	assert.Nil(t, err)

	queryDoc.Url, err = url.Parse("http://localhost:40123")
	assert.Nil(t, err)

	fetcher.scrapeLinks(queryDoc, pageResults)

	assert.Equal(t, expectedPageResults, pageResults)
}
