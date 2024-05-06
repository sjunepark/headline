package rodext

import (
	"github.com/cockroachdb/errors"
	"github.com/go-rod/rod"
	"github.com/sejunpark/headline/internal/pkg/util"
)

type Element struct {
	rodElement *rod.Element
}

func newElements(elements rod.Elements) []*Element {
	newEls := make([]*Element, len(elements))
	for i, el := range elements {
		newEls[i] = &Element{rodElement: el}
	}
	return newEls
}

func (e *Element) Element(selector string) (*Element, error) {
	elements, err := e.rodElement.Elements(selector)
	if err != nil {
		html, _ := e.rodElement.HTML()
		return nil, errors.Wrapf(err, "e.rodElement.Elements(%s) failed for html %s", selector, util.TrimLog(html))
	}
	if len(elements) > 1 {
		html, _ := e.rodElement.HTML()
		return nil, errors.Wrapf(MultipleElementsFoundError, "multiple elements found for selector %s in html %s", selector, util.TrimLog(html))
	}
	if len(elements) == 0 {
		html, _ := e.rodElement.HTML()
		return nil, errors.Wrapf(ElementNotFoundError, "element not found for selector %s in html %s", selector, util.TrimLog(html))
	}
	return &Element{rodElement: elements[0]}, nil
}

func (e *Element) Elements(selector string) ([]*Element, error) {
	elements, err := e.rodElement.Elements(selector)
	if err != nil {
		html, _ := e.rodElement.HTML()
		return nil, errors.Wrapf(err, "e.rodElement.Elements(%s) failed for html %s", selector, util.TrimLog(html))
	}

	if len(elements) == 0 {
		html, _ := e.rodElement.HTML()
		return nil, errors.Wrapf(ElementNotFoundError, "element not found for selector %s in html %s", selector, util.TrimLog(html))
	}

	return newElements(elements), nil
}

// Attribute returns the value of the attribute with the given name.
// If the attribute is not found, it returns an ElementNotFoundError.
// If the attribute is found but the value is empty, it returns an empty string.
func (e *Element) Attribute(name string) (string, error) {
	attribute, err := e.rodElement.Attribute(name)
	if err != nil {
		html, _ := e.rodElement.HTML()
		return "", errors.Wrapf(err, "e.rodElement.Attribute(%s) failed for html %s", name, util.TrimLog(html))
	}
	if attribute == nil {
		html, _ := e.rodElement.HTML()
		return "", errors.Wrapf(ElementNotFoundError, "attribute not found for name %s in html %s", name, util.TrimLog(html))
	}
	return *attribute, nil
}

// Text returns the text content of the element.
// If the text content is not found, an empty string is returned.
func (e *Element) Text() string {
	text, err := e.rodElement.Text()
	if err != nil {
		return ""
	}
	return text
}

// Equal returns true if the text content of the two elements' html content are the same.
// If one of them is nil, it returns a util.NilError.
func (e *Element) Equal(other *Element) (bool, error) {
	if e == nil || other == nil {
		return false, errors.Wrapf(util.NilError, "e.Equal(other), e: %v, other: %v", e, other)
	}

	eHTML, err := e.HTML()
	if err != nil {
		return false, errors.Wrapf(err, "e.HTML() failed for element %v", e)
	}
	otherHTML, err := other.HTML()
	if err != nil {
		return false, errors.Wrapf(err, "other.HTML() failed for element %v", other)
	}

	// You can't use rod's Equal method directly
	// because it's only available when two elements are in the same javascript world.
	equal := eHTML == otherHTML
	return equal, nil
}

func (e *Element) HTML() (string, error) {
	html, err := e.rodElement.HTML()
	if err != nil {
		return "", errors.Wrapf(err, "e.rodElement.HTML() failed for element %v", e)
	}
	return html, nil
}
