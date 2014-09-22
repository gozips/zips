package zips

import "archive/zip"
import "io"

type Zip struct {
	Sources []string
	errors  []error
}

func NewZip(srcs ...string) (z *Zip) {
	z = &Zip{
		Sources: srcs,
	}

	return
}

// Errors returns the errors collected during the zipping process
func (z Zip) Errors() []error {
	return z.errors
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

// FromFunc is the signature for the reader function that reads from the sources
// for each entry
type FromFunc func(string) (string, interface{})

// Readfrom calls f(s)
func (f FromFunc) Readfrom(s string) (string, interface{}) {
	return f(s)
}

// Write writes each file of the zip out to the Writer, passing each source
// through `srcfn` to get the name, content, error for each entry.
//
// `srcfn` can return either a ReadCloser or Error. This does allow one to
// return a ReadCloser on an error to create an entry containing that error
// message. Errors will skip creating an entry
func (z *Zip) Write(w io.Writer, fn FromFunc) (int64, bool) {
	var n int64
	ok := true
	zipOut := zip.NewWriter(w)
	defer zipOut.Close()

	for _, srcStr := range z.Sources {
		name, v := fn.Readfrom(srcStr)
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
