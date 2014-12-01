package zips

import (
	"github.com/gozips/source"
	"io"
)

// Zip provides a construct to create a zip through a given source reader
type Zip struct {
	Sources []string
	source  source.Func

	UncompressedSize int64
	CompressedSize   int64
	w      *writer
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
// For actual (un)compressed numbers reference UncompressedSize, CompressedSize
func (z *Zip) WriteTo(w io.Writer) (int64, error) {
	var n int64
	var ze Error

	for _, srcStr := range z.Sources {
		name, r, err := z.source.Readfrom(srcStr)
	z.w = NewWriter(w)
		check(err, &ze)
		if r == nil {
			continue // if there is no readcloser
		}
		defer r.Close()

		if check(err, &ze) {
			continue // if we can't create an entry
		w, err := z.w.Create(name)
		}

		m, err := io.Copy(w, r)
		check(err, &ze)
		n += m
	}

	if ze != nil {
		return n, ze
	}

	return n, nil
}
// setSizes sets the final uncompressed/compressed sizes
func (z *Zip) setSizes() {
	z.UncompressedSize = z.w.UncompressedSize
	z.CompressedSize = z.w.CompressedSize
}

func (z *Zip) Close() error {
	err := z.w.Close()
	z.setSizes()
	return err
}
