package crawler

import (
	"fmt"
	"io"
)

func EntryPoint(cliArgs []string, stdout, stderr io.Writer) int {
	if len(cliArgs) == 0 || cliArgs[0] == "" {
		fmt.Fprintf(stderr, "Error: please provide a url as the 1st argument\n")
		return 1
	}

	crawler := NewCrawlerApp(stdout, &HTTPFetcher{})
	if err := crawler.Run(cliArgs[0], stderr); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	return 0
}
