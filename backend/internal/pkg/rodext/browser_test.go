package rodext

import (
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
)

type BrowserSuite struct {
	BaseBrowserSuite
}

func (ts *BrowserSuite) SetupTest() {
	ts.SetupBaseBrowserTest()
}

func (ts *BrowserSuite) TearDownTest() {
	ts.TearDownBaseBrowserTest()
}

func TestBrowserSuite(t *testing.T) {
	suite.Run(t, new(BrowserSuite))
}

func (ts *BrowserSuite) TestBrowser_newBrowser() {
	ts.Run("all pages in pagePool should be nil after initialization", func() {
		size := ts.Browser.pagePool.len()
		for i := 0; i < size; i++ {
			p := <-ts.Browser.pagePool.pool
			ts.Nil(p)
		}
	})
}

func (ts *BrowserSuite) TestBrowser_Page() {
	multipliers := []int{1, 5, 100}

	ts.Run("unique pages returned from pool should be limited to the pool size", func() {
		for _, multiplier := range multipliers {
			createdPageIds := make(map[proto.TargetTargetID]bool)

			for i := 0; i < ts.PagePoolSize*multiplier; i++ {
				p, putPage, err := ts.Browser.Page()
				ts.NoError(err, "failed to get Page")

				id := p.rodPage.TargetID
				createdPageIds[id] = true
				putPage()
			}

			ts.Equal(len(createdPageIds), ts.PagePoolSize)
		}
	})

	ts.Run("when concurrent unique pages returned from pool should be limited to the pool size", func() {
		for _, multiplier := range multipliers {
			createdPageIds := make(chan proto.TargetTargetID, ts.PagePoolSize*multiplier)

			var wg sync.WaitGroup

			for i := 0; i < ts.PagePoolSize*multiplier; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					p, putPage, err := ts.Browser.Page()
					ts.NoError(err, "failed to get Page")
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
			ts.Equal(ts.PagePoolSize, len(pageIdsMap))
		}
	})
}
