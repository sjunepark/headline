package scraper

import (
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type PageSuite struct {
	suite.Suite
	browser       *Browser
	page          *Page
	pagePool      pagePool
	pagePoolLimit int
}

func (s *PageSuite) SetupTest() {
	browserOptions := BrowserOptions{
		NoDefaultDevice: true,
		Incognito:       true,
		Debug:           false,
	}
	b, err := NewBrowser(browserOptions)
	s.NoErrorf(err, "failed to initialize browser: %v", err)
	s.browser = b

	options := pageOptions{
		WindowFullscreen: true,
	}
	p, err := newPage(b, options)
	s.NoErrorf(err, "failed to initialize page: %v", err)
	s.page = p

	s.pagePoolLimit = 16
	s.pagePool = newPagePool(s.pagePoolLimit)
}

func (s *PageSuite) TearDownTest() {
	err := s.page.cleanup()
	s.NoErrorf(err, "failed to cleanup page: %v", err)

	err = s.page.rodPage.Browser().Close()
	s.NoErrorf(err, "failed to cleanup browser: %v", err)
}

func (s *PageSuite) TestPage_navigate() {
	err := s.page.navigate("https://www.wikipedia.org/")
	s.NoErrorf(err, "failed to navigate: %v", err)

	element, err := s.page.Element("h1")
	s.NoErrorf(err, "failed to find element: %v", err)

	text, err := element.Text()
	s.NoErrorf(err, "failed to get text: %v", err)
	s.True(strings.Contains(text, "Wikipedia"))
}

func (s *PageSuite) TestPagePool_get() {
	s.Run("unique pages returned from pool should be limited to the pool size", func() {
		create := newPageFactory(s.browser, pageOptions{WindowFullscreen: false})
		pageIds := make(map[proto.TargetTargetID]bool)

		multiplier := 2

		for i := 0; i < s.pagePoolLimit*multiplier; i++ {
			page, err := s.pagePool.get(create)
			s.NoErrorf(err, "failed to get page: %v", err)

			id := page.rodPage.TargetID
			s.T().Logf("got page: %v", id)

			pageIds[id] = true

			s.pagePool.put(page)
		}

		var numberOfUniqueIds int
		for i := 0; i < len(pageIds); i++ {
			numberOfUniqueIds++
		}

		s.Equal(s.pagePoolLimit, numberOfUniqueIds)
	})
}

func TestPageSuite(t *testing.T) {
	suite.Run(t, new(PageSuite))
}
