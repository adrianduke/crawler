package crawler

import (
	"errors"
	"fmt"
	"io"
	"net/url"
)

type CrawlerApp struct {
	Output  io.Writer
	Visited map[string]bool
	Fetcher
}

func NewCrawlerApp(output io.Writer, fetcher Fetcher) *CrawlerApp {
	return &CrawlerApp{
		Output:  output,
		Fetcher: fetcher,
		Visited: make(map[string]bool),
	}
}

func (ca *CrawlerApp) Run(rootURL *url.URL, depth int) error {
	if depth <= 0 {
		return errors.New("Reached max depth")
	}

	rootURL.Fragment = ""

	if _, found := ca.Visited[rootURL.String()]; found {
		return nil
	} else if !found {
		ca.Visited[rootURL.String()] = true
	}

	results, err := ca.Fetch(rootURL.String())
	if err != nil {
		return err
	}

	ca.PrettyPrint(rootURL, results)

	for internalURLString, _ := range results.internalURLs {
		internalURL, err := url.Parse(internalURLString)
		if err != nil {
			return err
		}

		if err := ca.Run(internalURL, depth-1); err != nil {
			return err
		}
	}

	return nil
}

func (ca *CrawlerApp) PrettyPrint(rootURL *url.URL, results *PageResults) {
	if rootURL.Path == "" {
		fmt.Fprintln(ca.Output, "/")
	} else {
		fmt.Fprintln(ca.Output, rootURL.Path)
	}

	fmt.Fprintln(ca.Output, "\tStatic Assets:")

	if len(results.cssStaticResources) > 0 {
		fmt.Fprintln(ca.Output, "\t\tCSS:")
		for _, cssResource := range results.cssStaticResources {
			fmt.Fprintf(ca.Output, "\t\t\t%s\n", cssResource)
		}
	}

	if len(results.jsStaticResources) > 0 {
		fmt.Fprintln(ca.Output, "\t\tJS:")
		for _, jsResource := range results.jsStaticResources {
			fmt.Fprintf(ca.Output, "\t\t\t%s\n", jsResource)
		}
	}

	fmt.Fprintln(ca.Output, "\tInternal Links:")
	if len(results.internalURLs) > 0 {
		for internalURL, _ := range results.internalURLs {
			fmt.Fprintf(ca.Output, "\t\t%s\n", internalURL)
		}
	}
	fmt.Fprintln(ca.Output)

}
