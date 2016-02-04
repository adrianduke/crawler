package crawler

import (
	"fmt"
	"io"
	"net/url"
)

func EntryPoint(cliArgs []string, stdout, stderr io.Writer) int {
	if len(cliArgs) == 0 || cliArgs[0] == "" {
		fmt.Fprintf(stderr, "Error: please provide a url as the 1st argument\n")
		return 1
	}

	targetURLString := cliArgs[0]

	targetURL, err := url.Parse(targetURLString)
	if err != nil {
		fmt.Fprintf(stderr, "Error: unable to parse url: %s\n", err)
		return 1
	}

	crawler := NewCrawlerApp(stdout, &HTTPFetcher{})

	err = crawler.Run(targetURL, 10)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %s\n", err)
		return 1
	}

	return 0
}
