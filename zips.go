package zips

import (
	"github.com/gozips/source"
	"io"
)

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

// WriteTo writes the zip out the Writer and returns the bytes that were *piped*
// through the zip.
// For actual (un)compressed numbers reference BytesIn, BytesOut
func (z *Zip) WriteTo(w io.Writer) (int64, error) {
	var n int64
	var ze Error

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

	err := zw.Close() // force close to calculate compression numbers
	check(err, &ze)
	z.BytesIn, z.BytesOut = zw.BytesIn, zw.BytesOut

	if ze != nil {
		return n, ze
	}
	return n, nil
}
