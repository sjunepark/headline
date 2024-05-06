package rodext

import (
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
		return nil, err
	}
	if len(elements) > 1 {
		return nil, MultipleElementsFoundError
	}
	if len(elements) == 0 {
		return nil, ElementNotFoundError
	}
	return &Element{rodElement: elements[0]}, nil
}

func (e *Element) Elements(selector string) ([]*Element, error) {
	elements, err := e.rodElement.Elements(selector)
	if err != nil {
		return nil, err
	}

	if len(elements) == 0 {
		return nil, ElementNotFoundError
	}

	return newElements(elements), nil
}

// Attribute returns the value of the attribute with the given name.
// If the attribute is not found, an empty string is returned.
func (e *Element) Attribute(name string) string {
	attribute, err := e.rodElement.Attribute(name)
	if err != nil {
		return ""
	}
	return *attribute
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

func (e *Element) Equal(other *Element) (bool, error) {
	if e == nil && other == nil {
		return true, nil
	}
	if e == nil || other == nil {
		return false, util.NilError
	}

	// You can't use rod's Equal method directly
	// because it's only available when two elements are in the same javascript world.
	equal := e.Text() == other.Text()
	return equal, nil
}
