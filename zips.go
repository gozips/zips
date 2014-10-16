package zips

import "fmt"
import "io"
import "strings"
import "github.com/gozips/source"

// Zip provides a trict around a zip.Writer
type Zip struct {
	Sources []string
	source  source.Func
}

// NewZip returns a zip that will read sources from the provided reader
func NewZip(fn source.Func) (z *Zip) {
	return &Zip{
		source: fn,
	}
}

// Add appends sources
func (z *Zip) Add(srcStr ...string) {
	z.Sources = append(z.Sources, srcStr...)
}

// check appends a ZipError
func check(e error, err *ZipError) bool {
	if e == nil {
		return false
	}

	*err = append(*err, e)
	return true
}

// ZipError is a collection of error that implements error
type ZipError []error

// Error returns a collective error
func (z ZipError) Error() string {
	var li []string
	for _, err := range z {
		li = append(li, fmt.Sprintf("* %s", err))
	}

	return fmt.Sprintf("%d error(s):\n\n%s", len(z), strings.Join(li, "\n"))
}

// WriteTo writes the zip out the Writer
func (z *Zip) WriteTo(w io.Writer) (int64, int64, error) {
	var ze ZipError
	var zw = NewWriter(w)
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

		_, err = io.Copy(w, r)
		check(err, &ze)
	}

	check(zw.Close(), &ze)
	return zw.BytesIn, zw.BytesOut, ze
}
