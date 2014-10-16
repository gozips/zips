package zips

import "fmt"
import "io"
import "strings"
import "github.com/gozips/source"

// Zip provides a construct to create a zip through a given source reader
type Zip struct {
	Sources []string
	source  source.Func

	BytesIn  int64
	BytesOut int64
}

func NewZip(fn source.Func) (z *Zip) {
	return &Zip{
		source: fn,
	}
}

// Add appends sources
func (z *Zip) Add(srcStr ...string) {
	z.Sources = append(z.Sources, srcStr...)
}

// ZipError is a collection of errors that implements error
type ZipError []error

// check appends a ZipError and returns a bool providing optional control flow
func check(e error, err *ZipError) bool {
	if e == nil {
		return false
	}

	*err = append(*err, e)
	return true
}

// Error returns a collective error
func (z ZipError) Error() string {
	var li []string
	for _, err := range z {
		li = append(li, fmt.Sprintf("* %s", err))
	}

	return fmt.Sprintf("%d error(s):\n\n%s", len(z), strings.Join(li, "\n"))
}

// WriteTo writes the zip out the Writer and returns the bytes that were *piped*
// through the zip.
// For actual (un)compressed numbers reference BytesIn, BytesOut
func (z *Zip) WriteTo(w io.Writer) (int64, error) {
	var n int64
	var ze ZipError

	zw := NewWriter(w)
	for _, srcStr := range z.Sources {
		name, r, err := z.source.Readfrom(srcStr)
		check(err, &ze)

		if r == nil {
			continue // if there is no readcloser
		}

		defer r.Close()
		w, err := zw.Create(name)
		if check(err, &ze) {
			continue // if we can't create an entry
		}

		m, err := io.Copy(w, r)
		check(err, &ze)
		n += m
	}

	check(zw.Close(), &ze)
	z.BytesIn, z.BytesOut = zw.BytesIn, zw.BytesOut

	return n, ze
}
