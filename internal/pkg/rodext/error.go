package rodext

import "github.com/cockroachdb/errors"

var MultipleElementsFoundError = errors.New("multiple elements found")

var ElementNotFoundError = errors.New("element not found")
