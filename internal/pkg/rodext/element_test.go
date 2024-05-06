package rodext

import (
	"github.com/sejunpark/headline/internal/pkg/util"
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
		aElement, err := ts.ulElement.Element("a[href='#section2']")
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

func (ts *ElementSuite) TestElement_Text() {
	tt := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "should return the text content of the element",
			got:  ts.ulElement.Text(),
			want: "Section 1\nSection 2\nSection 3\nSection 1",
		},
	}

	for _, tc := range tt {
		ts.Run(tc.name, func() {
			text := ts.ulElement.Text()
			ts.Equal(tc.want, text)
		})

	}
}

func (ts *ElementSuite) TestElement_HTML() {
	tt := []struct {
		name string
		want string
	}{
		{
			name: "should return the HTML content of the element",
			want: "<ul class=\"ulClass\" empty=\"\">\n        <li><a href=\"#section1\">Section 1</a></li>\n        <li><a href=\"#section2\">Section 2</a></li>\n        <li><a href=\"#section3\">Section 3</a></li>\n        <li><a href=\"#section1\">Section 1</a></li>\n    </ul>",
		},
	}

	for _, tc := range tt {
		ts.Run(tc.name, func() {
			html, err := ts.ulElement.HTML()
			ts.NoError(err, "failed to get HTML")
			ts.Equal(tc.want, html)
		})
	}
}

func (ts *ElementSuite) TestElement_Equal() {
	ts.Run("should return true if the text content of the two elements html are the same", func() {
		firstLi, err := ts.ulElement.Element("li:nth-child(1)")
		ts.NoError(err, "failed to get Element")

		fourthLi, err := ts.ulElement.Element("li:nth-child(4)")
		ts.NoError(err, "failed to get Element")

		ts.True(firstLi.Equal(firstLi))
		ts.True(firstLi.Equal(fourthLi))
		ts.True(fourthLi.Equal(firstLi))
	})

	ts.Run("should return false if the text content of the two elements html are not the same", func() {
		liSection2, err := ts.Page.Element("a[href='#section2']")
		ts.NoError(err, "failed to get Element")

		h3Section2, err := ts.Page.Element("h3#section2")
		ts.NoError(err, "failed to get Element")

		ts.False(liSection2.Equal(h3Section2))
		ts.False(h3Section2.Equal(liSection2))
	})

	ts.Run("should return false when both elements are nil", func() {
		var el1, el2 *Element
		ts.False(el1.Equal(el1))
		ts.False(el1.Equal(el2))
		ts.False(el2.Equal(el1))
	})

	ts.Run("should return false when one of the elements is nil", func() {
		var el *Element
		equal, err := ts.ulElement.Equal(el)
		ts.ErrorIs(err, util.NilError)
		ts.False(equal)

		equal, err = el.Equal(ts.ulElement)
		ts.ErrorIs(err, util.NilError)
		ts.False(equal)
	})
}

func isElement(toCheck any) bool {
	switch toCheck.(type) {
	case *Element:
		return true
	}
	return false
}
