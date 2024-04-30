package scraper

import (
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/suite"
	"strings"
	"sync"
	"testing"
)

type BrowserSuite struct {
	suite.Suite
	browser       *Browser
	pagePoolLimit int
}

func (ts *BrowserSuite) SetupTest() {
	browserOptions := BrowserOptions{
		NoDefaultDevice: true,
		Incognito:       true,
		Debug:           false,
		PagePoolSize:    16,
	}
	b, err := NewBrowser(browserOptions)
	ts.NoErrorf(err, "failed to initialize browser: %v", err)
	ts.browser = b

	ts.pagePoolLimit = browserOptions.PagePoolSize
}

func (ts *BrowserSuite) TearDownTest() {
	err := ts.browser.Cleanup()
	ts.NoErrorf(err, "failed to cleanup page: %v", err)
}

func (ts *BrowserSuite) TestPage_navigate() {
	page, putPage, err := ts.browser.Page()
	defer putPage()
	ts.NoErrorf(err, "failed to get page: %v", err)

	err = page.navigate("https://www.wikipedia.org/")
	ts.NoErrorf(err, "failed to navigate: %v", err)

	element, err := page.Element("h1")
	ts.NoErrorf(err, "failed to find element: %v", err)

	text, err := element.Text()
	ts.NoErrorf(err, "failed to get text: %v", err)
	ts.True(strings.Contains(text, "Wikipedia"))
}

func (ts *BrowserSuite) TestBrowser_Page() {
	multipliers := []int{1, 2, 100}

	ts.Run("unique pages returned from pool should be limited to the pool size", func() {
		for _, multiplier := range multipliers {
			pageIds := make(map[proto.TargetTargetID]bool)

			for i := 0; i < ts.pagePoolLimit*multiplier; i++ {
				page, putPage, err := ts.browser.Page()
				ts.NoErrorf(err, "failed to get page: %v", err)

				id := page.rodPage.TargetID
				pageIds[id] = true
				putPage()
			}

			ts.Equal(ts.pagePoolLimit, len(pageIds))
		}

	})

	ts.Run("works on concurrent access", func() {
		for _, multiplier := range multipliers {
			pageIds := make(chan proto.TargetTargetID, ts.pagePoolLimit*multiplier)

			var wg sync.WaitGroup
			for i := 0; i < ts.pagePoolLimit*multiplier; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					page, putPage, err := ts.browser.Page()
					ts.NoErrorf(err, "failed to get page: %v", err)

					id := page.rodPage.TargetID
					pageIds <- id
					putPage()
				}()
			}
			wg.Wait()
			close(pageIds)

			pageIdsMap := make(map[proto.TargetTargetID]bool)
			for id := range pageIds {
				pageIdsMap[id] = true
			}
			ts.Equal(ts.pagePoolLimit, len(pageIdsMap))
		}
	})
}

func TestPageSuite(t *testing.T) {
	suite.Run(t, new(BrowserSuite))
}
