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

// Write writes each file of the zip out to the Writer, passing each source
// through `srcfn` to get the name, content, error for each entry.
//
// If `srcfn` returns only an Error, the entry will be skipped. If it returns
// an Error and a ReadCloser it will stack the Error and create an entry
// with the content of the ReadCloser. This allows one to create entries for
// unprocessible entries with the reason as the content of the entry file.
// When writing handling an error as an entry, entry names should be appended
// with .txt for easy opening and reading
func (z *Zip) Write(srcfn func(string) (string, io.ReadCloser, error), w io.Writer) (int64, bool) {
	var n int64
	ok := true
	zipOut := zip.NewWriter(w)
	defer zipOut.Close()

	for _, srcStr := range z.Sources {
		name, r, err := srcfn(srcStr)
		if z.check(err, &ok) && r == nil {
			continue // only if err && ReadCloser is nil
		}
		defer r.Close()

		w, err := zipOut.Create(name)
		if z.check(err, &ok) {
			continue // if we can't create an entry
		}

		m, err := io.Copy(w, r)
		z.check(err, &ok)

		n += m
	}

	return n, ok
}
