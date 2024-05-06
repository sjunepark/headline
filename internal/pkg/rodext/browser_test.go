package rodext

import (
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

type BrowserSuite struct {
	suite.Suite
	browser        *Browser
	cleanupBrowser func()
	pagePoolSize   int
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
	ts.pagePoolSize = options.PagePoolSize
}

func (ts *BrowserSuite) TearDownTest() {
	ts.cleanupBrowser()
	ts.Equal(ts.browser.pagePool.len(), 0, "pagePool should be empty after Cleanup")
}

func TestBrowserSuite(t *testing.T) {
	suite.Run(t, new(BrowserSuite))
}

func (ts *BrowserSuite) TestBrowser_newBrowser() {
	ts.Run("all pages in pagePool should be nil after initialization", func() {
		size := ts.browser.pagePool.len()
		for i := 0; i < size; i++ {
			p := <-ts.browser.pagePool.pool
			ts.Nil(p)
		}
	})
}

func (ts *BrowserSuite) TestBrowser_Page() {
	multipliers := []int{1, 5, 100}

	ts.Run("unique pages returned from pool should be limited to the pool size", func() {
		for _, multiplier := range multipliers {
			createdPageIds := make(map[proto.TargetTargetID]bool)

			for i := 0; i < ts.pagePoolSize*multiplier; i++ {
				p, putPage, err := ts.browser.Page()
				ts.NoErrorf(err, "failed to get Page: %v", err)

				id := p.rodPage.TargetID
				createdPageIds[id] = true
				putPage()
			}

			ts.Equal(len(createdPageIds), ts.pagePoolSize)
		}
	})

	ts.Run("concurrent: unique pages returned from pool should be limited to the pool size", func() {
		for _, multiplier := range multipliers {
			createdPageIds := make(chan proto.TargetTargetID, ts.pagePoolSize*multiplier)

			var wg sync.WaitGroup

			for i := 0; i < ts.pagePoolSize*multiplier; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					p, putPage, err := ts.browser.Page()
					ts.NoErrorf(err, "failed to get Page: %v", err)
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
			ts.Equal(ts.pagePoolSize, len(pageIdsMap))
		}
	})
}
