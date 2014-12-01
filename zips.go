package zips

import (
	"errors"
	"github.com/gozips/source"
	"io"
)

// Zip provides a construct to create a zip through a given source reader
type Zip struct {
	Sources          []string
	UncompressedSize int64
	CompressedSize   int64
	N                int64 // sum of total bytes written through each Entry call

	source source.Func
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
	var e Error

	z.w = NewWriter(w)
	for _, v := range z.Sources {
		name, r, err := z.source.Readfrom(v)
		check(err, &e)
		if r == nil {
			continue // if there is no readcloser
		}
		defer r.Close()

		m, err := z.AddEntry(name, r)
		check(err, &e)
		n += m
	}

	if e != nil {
		return n, e
	}

	return n, nil
}

// incN increments N by n and returns the incremented N
func (z *Zip) incN(n int64) int64 {
	z.N += n
	return z.N
}

// AddEntry creates a new zip entry. This function is not intented for general
// purpose use and as such WriteTo is required to be called before AddEntry can
// be used.
// Any call to AddEntry will increment N
func (z *Zip) AddEntry(name string, r io.Reader) (int64, error) {
	if z.w == nil {
		return 0, errors.New("error: writer: undefined")
	}

	w, err := z.w.Create(name)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(w, r)
	z.incN(n)

	return n, err
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
