package crawler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type PageResults struct {
	jsStaticResources  []string
	cssStaticResources []string
	imgStaticResources map[string]bool
	internalURLs       map[string]bool
	externalURLs       map[string]bool
	errors             []error
}

func NewPageResults() *PageResults {
	return &PageResults{
		jsStaticResources:  make([]string, 0),
		cssStaticResources: make([]string, 0),
		imgStaticResources: make(map[string]bool),
		internalURLs:       make(map[string]bool),
		externalURLs:       make(map[string]bool),
		errors:             make([]error, 0),
	}
}

type Fetcher interface {
	Fetch(url string) (pageResults *PageResults, err error)
}

type FetcherAdapter func(string) (*PageResults, error)

func (f FetcherAdapter) Fetch(url string) (*PageResults, error) {
	return f(url)
}

type HTTPFetcher struct{}

func (hf *HTTPFetcher) Fetch(url string) (*PageResults, error) {
	results := NewPageResults()

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch url: %s", err)
	}
	defer resp.Body.Close()

	queryDoc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	hf.scrapeStaticResources(queryDoc, results)
	hf.scrapeLinks(queryDoc, results)

	return results, nil
}

func (hf *HTTPFetcher) scrapeStaticResources(queryDoc *goquery.Document, pageResults *PageResults) {
	queryDoc.Find("link").Each(func(i int, s *goquery.Selection) {
		if href := s.AttrOr("href", ""); href != "" {
			parsedHref, err := url.Parse(href)
			if err != nil {
				pageResults.errors = append(pageResults.errors, err)
				return
			}
			if !parsedHref.IsAbs() || parsedHref.Host == queryDoc.Url.Host {
				parsedHref = fixupURL(parsedHref, queryDoc.Url)
			}

			pageResults.cssStaticResources = append(pageResults.cssStaticResources, parsedHref.String())
		}
	})

	queryDoc.Find("script").Each(func(i int, s *goquery.Selection) {
		if src := s.AttrOr("src", ""); src != "" {
			parsedSrc, err := url.Parse(src)
			if err != nil {
				pageResults.errors = append(pageResults.errors, err)
				return
			}
			if !parsedSrc.IsAbs() || parsedSrc.Host == queryDoc.Url.Host {
				parsedSrc = fixupURL(parsedSrc, queryDoc.Url)
			}

			pageResults.jsStaticResources = append(pageResults.jsStaticResources, parsedSrc.String())
		}
	})

	queryDoc.Find("img").Each(func(i int, s *goquery.Selection) {
		if src := s.AttrOr("src", ""); src != "" {
			parsedSrc, err := url.Parse(src)
			if err != nil {
				pageResults.errors = append(pageResults.errors, err)
				return
			}

			if !parsedSrc.IsAbs() || parsedSrc.Host == queryDoc.Url.Host {
				parsedSrc = fixupURL(parsedSrc, queryDoc.Url)
			}

			if _, found := pageResults.imgStaticResources[parsedSrc.String()]; !found {
				pageResults.imgStaticResources[parsedSrc.String()] = true
			}
		}
	})
}

func (hf *HTTPFetcher) scrapeLinks(queryDoc *goquery.Document, pageResults *PageResults) {
	queryDoc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href := s.AttrOr("href", ""); href != "" {
			parsedHref, err := url.Parse(href)
			if err != nil {
				pageResults.errors = append(pageResults.errors, err)
				return
			}

			if !parsedHref.IsAbs() || parsedHref.Host == queryDoc.Url.Host {
				parsedHref = fixupURL(parsedHref, queryDoc.Url)
				if _, found := pageResults.internalURLs[parsedHref.String()]; !found {
					pageResults.internalURLs[parsedHref.String()] = true
				}
			} else {
				if _, found := pageResults.internalURLs[parsedHref.String()]; !found {
					pageResults.externalURLs[parsedHref.String()] = true
				}
			}
		}
	})
}

// Resolve relative urls to absolute, strip fragments
func fixupURL(incoming, base *url.URL) *url.URL {
	absoluteURL := base.ResolveReference(incoming)
	absoluteURL.Fragment = ""

	return absoluteURL
}
