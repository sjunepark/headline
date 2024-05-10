package rodext

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type PageSuite struct {
	BasePageSuite
}

func (ts *PageSuite) SetupSuite() {
	ts.SetupBasePageSuite()
}

func (ts *PageSuite) TearDownSuite() {
	ts.TearDownBasePageSuite()
}

func (ts *PageSuite) SetupSubTest() {
	ts.SetupBasePageSubTest()
}

func (ts *PageSuite) TearDownSubTest() {
	ts.TearDownBasePageSubTest()
}

func TestPageSuite(t *testing.T) {
	suite.Run(t, new(PageSuite))
}

func (ts *PageSuite) TestPage_cleanup() {
	ts.Run("page should be cleaned up after cleanup is called", func() {
		page := ts.Page
		_, err := page.rodPage.Info()
		ts.NoError(err, "failed to get Page info")

		page.cleanup()

		_, err = page.rodPage.Info()
		ts.Error(err, "shouldn't be able to get Page info after Cleanup")
	})
}

func (ts *PageSuite) TestPage_Element() {
	ts.Run("Element should return the element if it exists", func() {
		el, err := ts.Page.Element("h1")
		ts.NoError(err, "failed to get element")
		ts.NotNil(el, "element should not be nil")
	})

	ts.Run("Element should return an error when multiple elements are found", func() {
		el, err := ts.Page.Element("a")
		ts.ErrorIs(err, MultipleElementsFoundError, "should return MultipleElementsFoundError when multiple elements are found")
		ts.Nil(el, "when multiple elements are found, element should be nil")
	})

	ts.Run("Element should return an error when no elements are found", func() {
		el, err := ts.Page.Element(".nonexistent-element")
		ts.ErrorIs(err, ElementNotFoundError, "should return ElementNotFoundError when no elements are found")
		ts.Nil(el, "when no elements are found, element should be nil")
	})
}

func (ts *PageSuite) TestPage_Elements() {
	ts.Run("Elements should return the elements if they exist", func() {
		els, err := ts.Page.Elements("a")
		ts.NoError(err, "failed to get elements")
		ts.NotEmpty(els, "elements should not be empty")
	})

	ts.Run("Elements should return an empty slice when no elements are found", func() {
		els, err := ts.Page.Elements(".nonexistent-element")
		ts.NoError(err, "failed to get elements")
		ts.Empty(els, "elements should be empty")
	})
}
