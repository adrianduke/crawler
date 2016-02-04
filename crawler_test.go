package crawler

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ItPrintsForwardSlashForURLsWithNoPath(t *testing.T) {
	output := new(bytes.Buffer)
	mockFetcher := &MockFetcher{}
	app := NewCrawlerApp(output, mockFetcher)

	url, err := url.Parse("http://www.google.com")
	assert.Nil(t, err)

	expectedOutput := `/
	Static Assets:
	Internal Links:

`

	mockFetcher.On("Fetch", url.String()).Return(&PageResults{}, nil)

	app.Run(url, 2)

	assert.Equal(t, output.String(), expectedOutput)
}

func Test_ItReturnsAnErrorOnceMaxDepthIsReached(t *testing.T) {
	output := new(bytes.Buffer)
	mockFetcher := &MockFetcher{}
	app := NewCrawlerApp(output, mockFetcher)

	url, err := url.Parse("http://www.google.com")
	assert.Nil(t, err)

	mockFetcher.On("Fetch", url.String()).Return(&PageResults{
		internalURLs: map[string]bool{"http://www.google.com/1": true},
	}, nil)

	err = app.Run(url, 1)

	assert.EqualError(t, err, "Reached max depth")
}

type MockFetcher struct {
	mock.Mock
}

func (m *MockFetcher) Fetch(url string) (*PageResults, error) {
	args := m.Called(url)

	return (args.Get(0)).(*PageResults), args.Error(1)
}