package zips

import "archive/zip"
import "io"
import "github.com/gozips/source"

// Zip provides a trict around a zip.Writer
type Zip struct {
	Sources []string
	errors  []error
	r       source.Func
}

// NewZip returns a zip that will read sources from the provided reader
func NewZip(r source.Func) (z *Zip) {
	return &Zip{
		r: r,
	}
}

// Errors returns the errors collected during the zipping process
func (z Zip) Errors() []error {
	return z.errors
}

// Add appends sources
func (z *Zip) Add(srcStr string) {
	z.Sources = append(z.Sources, srcStr)
}

// check adds *unhandlable* errors to a slice to be later inspected.
// Also marks `ok` as false
func (z *Zip) check(e error, ok *bool) bool {
	if e == nil {
		return false
	}

	z.errors = append(z.errors, e)
	*ok = false

	return true
}

// WriteTo writes the zip out the Writer
func (z *Zip) WriteTo(w io.Writer) (int64, bool) {
	var n int64
	ok := true
	zipOut := zip.NewWriter(w)
	defer zipOut.Close()

	for _, srcStr := range z.Sources {
		name, v := z.r.Readfrom(srcStr)
		switch r := v.(type) {
		case io.ReadCloser:
			defer r.Close()

			w, err := zipOut.Create(name)
			if z.check(err, &ok) {
				continue // if we can't create an entry
			}

			m, err := io.Copy(w, r)
			z.check(err, &ok)

			n += m
		case error:
			z.check(r, &ok)
		}
	}

	return n, ok
}
