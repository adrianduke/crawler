package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DATA-DOG/godog"
	gherkin "github.com/cucumber/gherkin-go"
)

const TempPrefix = "crawler"

func NewWebCrawlingFeature() *webCrawlingFeature {
	return &webCrawlingFeature{
		Stdout:   new(bytes.Buffer),
		Stderr:   new(bytes.Buffer),
		Cleanups: make([]func() error, 0),
	}
}

type webCrawlingFeature struct {
	Stdout     *bytes.Buffer
	Stderr     *bytes.Buffer
	ExitCode   int
	Cleanups   []func() error
	RootDir    string
	TestServer *httptest.Server
}

func (wcf *webCrawlingFeature) Reset(v interface{}) {
	wcf.Stdout = new(bytes.Buffer)
	wcf.Stderr = new(bytes.Buffer)
	wcf.ExitCode = 0
	wcf.RootDir = ""

	for _, cleanupFn := range wcf.Cleanups {
		err := cleanupFn()
		if err != nil {
			fmt.Println("Cleanup Error:", err.Error())
		}
	}
	wcf.Cleanups = make([]func() error, 0)
	wcf.TestServer = nil
}

func (wcf *webCrawlingFeature) Run(cliArgs []string) {
	wcf.ExitCode = EntryPoint(cliArgs, wcf.Stdout, wcf.Stderr)
}

func featureContext(s *godog.Suite) {
	wcf := NewWebCrawlingFeature()

	s.BeforeScenario(wcf.Reset)

	s.Step(`^an invalid url "([^"]*)"$`, wcf.anInvalidUrl)
	s.Step(`^I run the crawler with the following arguments "([^"]*)"$`, wcf.iRunTheCrawlerWithTheFollowingArguments)
	s.Step(`^I should see an error informing me "([^"]*)"$`, wcf.iShouldSeeAnErrorInformingMe)
	s.Step(`^a unfetchable url "([^"]*)"$`, wcf.aUnfetchableUrl)
	s.Step(`^the exit code should be: (\d+)$`, wcf.theExitCodeShouldBe)
	s.Step(`^a webpage "([^"]*)" containing:$`, wcf.aWebpagecontaining)
	s.Step(`^the webpages are being hosted locally$`, wcf.theWebpagesAreBeingHostedLocally)
	s.Step(`^I run the crawler with the locally hosted url$`, wcf.iRunTheCrawlerWithTheLocallyHostedUrl)
	s.Step(`^I should see the following:$`, wcf.iShouldSeeTheFollowing)
}

func (wcf *webCrawlingFeature) anInvalidUrl(arg1 string) error {
	return nil
}

func (wcf *webCrawlingFeature) iRunTheCrawlerWithTheFollowingArguments(cliArgs string) error {
	wcf.Run(strings.Split(cliArgs, " "))

	return nil
}

func (wcf *webCrawlingFeature) iShouldSeeAnErrorInformingMe(errorRegex string) error {
	matched, err := regexp.MatchString(errorRegex, wcf.Stderr.String())
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("Unable to find match for '%s' in:\n\n%s", errorRegex, wcf.Stderr.String())
	}

	return nil
}

func (wcf *webCrawlingFeature) aUnfetchableUrl(arg1 string) error {
	return nil
}

func (wcf *webCrawlingFeature) theExitCodeShouldBe(expectedExitCode int) error {
	if wcf.ExitCode != expectedExitCode {
		return fmt.Errorf("Exit code '%d' did not match expected '%d'", wcf.ExitCode, expectedExitCode)
	}

	return nil
}

func (wcf *webCrawlingFeature) aWebpagecontaining(relativeFilePath string, contents *gherkin.DocString) error {
	var err error

	if wcf.RootDir == "" {
		wcf.RootDir, err = ioutil.TempDir("", TempPrefix)
		if err != nil {
			return err
		}
		wcf.Cleanups = append(wcf.Cleanups, func() error { return os.RemoveAll(wcf.RootDir) })
	}

	absFilePath := filepath.Join(wcf.RootDir, relativeFilePath)
	err = ioutil.WriteFile(absFilePath, []byte(contents.Content), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (wcf *webCrawlingFeature) theWebpagesAreBeingHostedLocally() error {
	wcf.TestServer = httptest.NewServer(http.FileServer(http.Dir(wcf.RootDir)))

	wcf.Cleanups = append(wcf.Cleanups, func() error { wcf.TestServer.Close(); return nil })

	return nil
}

func (wcf *webCrawlingFeature) iRunTheCrawlerWithTheLocallyHostedUrl() error {
	if wcf.TestServer == nil {
		return errors.New("TestServer needs to be setup before it can be crawled")
	}

	wcf.Run([]string{wcf.TestServer.URL})

	return nil
}

func (wcf *webCrawlingFeature) iShouldSeeTheFollowing(outputRegex *gherkin.DocString) error {
	matched, err := regexp.MatchString(outputRegex.Content, wcf.Stdout.String())
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("Unable to find match for '%s' in:\n\n%s", outputRegex.Content, wcf.Stdout.String())
	}

	return nil
}
