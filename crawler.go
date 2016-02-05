package crawler

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"sync"
)

type CrawlerApp struct {
	Output    io.Writer
	Visited   map[string]bool
	Errors    chan error
	waitGroup *sync.WaitGroup
	mutex     *sync.Mutex
	Fetcher
}

func NewCrawlerApp(output io.Writer, fetcher Fetcher) *CrawlerApp {
	return &CrawlerApp{
		Output:    output,
		Fetcher:   fetcher,
		Errors:    make(chan error),
		Visited:   make(map[string]bool),
		waitGroup: &sync.WaitGroup{},
		mutex:     &sync.Mutex{},
	}
}

func (ca *CrawlerApp) Run(rootURLString string, errorOutput io.Writer) error {
	rootURL, err := url.Parse(rootURLString)
	if err != nil {
		return fmt.Errorf("Error: unable to parse url: %s\n", err)
	}

	doneCh := make(chan struct{})
	go func() {
		ca.waitGroup.Wait()
		close(doneCh)
	}()

	ca.waitGroup.Add(1)
	go ca.Crawl(rootURL, 5)

	var hasErrored bool
	for {
		select {
		case err := <-ca.Errors:
			fmt.Fprintf(errorOutput, "Error: %s\n", err)
			hasErrored = true
		case <-doneCh:
			if hasErrored {
				return errors.New("")
			}

			return nil
		}
	}
}

func (ca *CrawlerApp) Crawl(rootURL *url.URL, depth int) {
	defer ca.waitGroup.Done()

	if depth <= 0 {
		ca.Errors <- errors.New("Reached max depth")
		return
	}

	rootURL.Fragment = ""

	ca.mutex.Lock()
	if _, found := ca.Visited[rootURL.String()]; found {
		ca.mutex.Unlock()
		return
	} else if !found {
		ca.Visited[rootURL.String()] = true
		ca.mutex.Unlock()
	}

	results, err := ca.Fetch(rootURL.String())
	if err != nil {
		ca.Errors <- err
		return
	}

	ca.PrettyPrint(rootURL, results)

	for internalURLString, _ := range results.internalURLs {
		internalURL, err := url.Parse(internalURLString)
		if err != nil {
			ca.Errors <- err
			return
		}

		ca.waitGroup.Add(1)
		go ca.Crawl(internalURL, depth-1)
	}
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
