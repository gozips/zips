package zips

import (
	"fmt"
	"strings"
)

// Error is a collection of errors that implements error
type Error []error

// check appends a Error and returns a bool providing optional control flow
func check(e error, err *Error) bool {
	if e == nil {
		return false
	}

	*err = append(*err, e)
	return true
}

// Error returns a collective error
func (z Error) Error() string {
	var li []string
	for _, err := range z {
		li = append(li, fmt.Sprintf("* %s", err))
	}

	return fmt.Sprintf("%d error(s):\n\n%s", len(z), strings.Join(li, "\n"))
}
