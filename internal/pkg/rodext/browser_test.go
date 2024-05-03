package rodext

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
	browser        *Browser
	cleanupBrowser func()
	pagePoolLimit  int
}

func (ts *BrowserSuite) SetupTest() {
	options := BrowserOptions{
		NoDefaultDevice: true,
		Incognito:       true,
		Debug:           false,
		PagePoolSize:    16,
	}
	b, cleanup, err := NewBrowser(options)
	ts.NoErrorf(err, "failed to initialize Browser: %v", err)
	ts.browser = b
	ts.cleanupBrowser = cleanup
	ts.pagePoolLimit = options.PagePoolSize
}

func (ts *BrowserSuite) TearDownTest() {
	ts.cleanupBrowser()
	ts.Equal(len(ts.browser.pagePool), 0, "pagePool should be empty after Cleanup")
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

func (ts *BrowserSuite) TestPage_cleanup() {
	ts.Run("page should be cleaned up after putPage is called", func() {
		p, cleanup, err := ts.browser.newPage(pageOptions{})
		ts.NoErrorf(err, "failed to create new page: %v", err)

		_, err = p.rodPage.Info()
		ts.NoErrorf(err, "failed to get page info: %v", err)

		cleanup()

		_, err = p.rodPage.Info()
		ts.Error(err, "shouldn't be able to get page info after Cleanup")
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

type BrowserCleanupSuite struct {
	suite.Suite
	browser        *Browser
	cleanupBrowser func()
	pagePoolLimit  int
}

func (ts *BrowserCleanupSuite) SetupTest() {
	options := BrowserOptions{
		NoDefaultDevice: true,
		Incognito:       true,
		Debug:           false,
		PagePoolSize:    16,
	}
	b, cleanup, err := NewBrowser(options)
	ts.NoErrorf(err, "failed to initialize Browser: %v", err)
	ts.browser = b
	ts.cleanupBrowser = cleanup
	ts.pagePoolLimit = options.PagePoolSize
}

func TestBrowserCleanupSuite(t *testing.T) {
	suite.Run(t, new(BrowserCleanupSuite))
}

func (ts *BrowserCleanupSuite) TestBrowser_cleanup() {
	ts.Run("pagePool should be empty after Cleanup", func() {
		ts.cleanupBrowser()
		ts.Equal(len(ts.browser.pagePool), 0, "pagePool should be empty after Cleanup")
	})
}

func (ts *BrowserCleanupSuite) TestBrowser_cancel() {
	ts.Run("on context cancellation, everything should be cleaned up gracefully", func() {
		//	todo: implement
	})
}
