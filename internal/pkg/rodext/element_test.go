package rodext

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ElementSuite struct {
	BaseElementSuite
}

func (ts *ElementSuite) SetupSuite() {
	ts.SetupBaseElementSuite()
}

func (ts *ElementSuite) TearDownSuite() {
	ts.TearDownBaseElementSuite()
}

func TestElementSuite(t *testing.T) {
	suite.Run(t, new(ElementSuite))
}

func (ts *ElementSuite) TestElement_Element() {
	ts.Run("should return an Element with the correct selector", func() {
		aElement, err := ts.ulElement.Element("a[href='#section1']")
		html, _ := aElement.HTML()
		ts.NoErrorf(err, "failed to get Element, target ulElement: %v", html)
		ts.Truef(isElement(aElement), "returned Element is not an Element, got %T", aElement)
		ts.NotNilf(aElement, "returned Element is nil")
	})

	ts.Run("should error if the selector is invalid", func() {
		el, err := ts.ulElement.Element("invalid")
		ts.ErrorIsf(err, ElementNotFoundError, "expected ElementNotFoundError but got %v", err)
		ts.Nilf(el, "expected Element to be nil but got %v", el)
	})

	ts.Run("should error when multiple elements are found", func() {
		el, err := ts.ulElement.Element("li")
		ts.ErrorIsf(err, MultipleElementsFoundError, "expected MultipleElementsFoundError but got %v", err)
		ts.Nilf(el, "expected Element to be nil but got %v", el)
	})
}

func (ts *ElementSuite) TestElement_Elements() {
	ts.Run("should return a slice of Elements with the correct selector", func() {
		liElements, err := ts.ulElement.Elements("li")
		ts.NoError(err, "failed to get Elements")
		ts.Truef(len(liElements) > 0, "expected slice of Elements to have length > 0 but got %v", liElements)
		for _, el := range liElements {
			ts.Truef(isElement(el), "returned Element is not an Element, got %T", el)
		}
	})

	ts.Run("should error if the selector is invalid", func() {
		els, err := ts.ulElement.Elements("invalid")
		ts.ErrorIs(err, ElementNotFoundError)
		ts.Nilf(els, "expected Elements to be nil but got %v", els)
	})
}

func (ts *ElementSuite) TestElement_Attribute() {
	ts.Run("should return an attribute value", func() {
		ulClass, err := ts.ulElement.Attribute("class")
		ts.NoError(err, "failed to get attribute")
		ts.Equal("ulClass", ulClass)
	})

	ts.Run("should error attribute is not found", func() {
		got, err := ts.ulElement.Attribute("invalid")
		ts.ErrorIs(err, ElementNotFoundError)
		ts.Empty(got)
	})

	ts.Run("should return an empty string if the attribute exists but is empty", func() {
		emptyAttr, err := ts.ulElement.Attribute("empty")
		ts.NoError(err, "failed to get attribute")
		ts.Equal("", emptyAttr)
	})
}

func isElement(toCheck any) bool {
	switch toCheck.(type) {
	case *Element:
		return true
	}
	return false
}
