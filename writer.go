package zips

import "archive/zip"
import "io"

type fileheaders []*zip.FileHeader

// total returns the totals from the fileheaders for total read (uncompressed)
// and total out (compressed)
func (f fileheaders) total() (int64, int64) {
	var u, c int64
	for _, fh := range f {
		u += int64(fh.UncompressedSize64)
		c += int64(fh.CompressedSize64)
	}

	return u, c
}

// Writer is a wrapper around zip.Writer to expose the uncompressed and
// compressed totals
type Writer struct {
	BytesIn  int64 // bytes read
	BytesOut int64 // bytes compressed out

	*zip.Writer
	fileHeaders fileheaders
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer: zip.NewWriter(w),
	}
}

// Create decomposes the original Create method and appends the FileHeader
func (z *Writer) Create(name string) (io.Writer, error) {
	fh := &zip.FileHeader{
		Name:   name,
		Method: zip.Deflate,
	}
	z.fileHeaders = append(z.fileHeaders, fh)

	return z.Writer.CreateHeader(fh)
}

// Close wraps the original close and calls tally
func (z *Writer) Close() error {
	err := z.Writer.Close()
	z.tally()

	return err
}

func (z *Writer) tally() {
	z.BytesIn, z.BytesOut = z.fileHeaders.total()
}
