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
		return nil, NotFoundError
	}
	return &Element{rodElement: elements[0]}, nil
}

func (e *Element) Elements(selector string) ([]*Element, error) {
	elements, err := e.rodElement.Elements(selector)
	if err != nil {
		return nil, err
	}

	if len(elements) == 0 {
		return nil, NotFoundError
	}

	return newElements(elements), nil
}

// Attribute returns the value of the attribute with the given name.
// If the attribute is not found, it returns a NotFoundError.
// If the attribute is found but the value is empty, it returns an empty string.
func (e *Element) Attribute(name string) (string, error) {
	attribute, err := e.rodElement.Attribute(name)
	if err != nil {
		return "", err
	}
	if attribute == nil {
		return "", NotFoundError
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

// Equal returns true if the text content of the two elements is the same.
// It returns true even when both elements are nil.
// If one of them is nil, it returns a util.NilError.
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
