package scraper

import (
	"fmt"
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/suite"
	"strings"
	"sync"
	"testing"
)

type BrowserSuite struct {
	suite.Suite
	browser        *browser
	cleanupBrowser func()
	pagePoolLimit  int
}

func (ts *BrowserSuite) SetupTest() {
	options := browserOptions{
		noDefaultDevice: true,
		incognito:       true,
		debug:           false,
		pagePoolSize:    16,
	}
	b, cleanup, err := newBrowser(options)
	ts.NoErrorf(err, "failed to initialize browser: %v", err)
	ts.browser = b
	ts.cleanupBrowser = cleanup
	ts.pagePoolLimit = options.pagePoolSize
}

func (ts *BrowserSuite) TearDownTest() {
	ts.cleanupBrowser()
}

func TestBrowserSuite(t *testing.T) {
	suite.Run(t, new(BrowserSuite))
}

func (ts *BrowserSuite) TestBrowser_newBrowser() {
	ts.Run("all pages in pagePool should be nil after initialization", func() {
		size := len(ts.browser.pagePool)
		for i := 0; i < size; i++ {
			p := <-ts.browser.pagePool
			ts.Nil(p)
		}
	})
}

func (ts *BrowserSuite) TestBrowser_Page() {
	multipliers := []int{1, 5, 100}

	ts.Run("unique pages returned from pool should be limited to the pool size", func() {
		for _, multiplier := range multipliers {
			createdPageIds := make(map[proto.TargetTargetID]bool)

			for i := 0; i < ts.pagePoolLimit*multiplier; i++ {
				p, putPage, err := ts.browser.page()
				ts.NoErrorf(err, "failed to get page: %v", err)

				id := p.rodPage.TargetID
				createdPageIds[id] = true
				putPage()
			}

			ts.Equal(len(createdPageIds), ts.pagePoolLimit)
		}

	})

	ts.Run("works on concurrent access", func() {
		for _, multiplier := range multipliers {
			createdPageIds := make(chan proto.TargetTargetID, ts.pagePoolLimit*multiplier)

			var wg sync.WaitGroup

			for i := 0; i < ts.pagePoolLimit*multiplier; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					fmt.Printf("running goroutine\n")
					p, putPage, err := ts.browser.page()
					ts.NoErrorf(err, "failed to get page: %v", err)
					if err == nil {
						defer putPage()
					}

					id := p.rodPage.TargetID
					createdPageIds <- id
				}()
			}
			wg.Wait()
			close(createdPageIds)

			pageIdsMap := make(map[proto.TargetTargetID]bool)
			for id := range createdPageIds {
				pageIdsMap[id] = true
			}
			ts.Equal(ts.pagePoolLimit, len(pageIdsMap))
		}
	})
}

func (ts *BrowserSuite) TestPage_navigate() {
	p, putPage, err := ts.browser.page()
	defer putPage()
	ts.NoErrorf(err, "failed to get page: %v", err)

	err = p.navigate("https://www.wikipedia.org/")
	ts.NoErrorf(err, "failed to navigate: %v", err)

	element, err := p.element("h1")
	ts.NoErrorf(err, "failed to find element: %v", err)

	text, err := element.Text()
	ts.NoErrorf(err, "failed to get text: %v", err)
	ts.True(strings.Contains(text, "Wikipedia"))
}
