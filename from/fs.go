package from

import "fmt"
import "io"
import "os"

// FS returns the file name, content for the given pathStr
// PathErrors will not be returned as an Error, instead it will return the error
// as content creating an entry for the errored file open
func FS(pathStr string) (string, io.ReadCloser, error) {
	var r io.ReadCloser
	var err error
	name := filename(pathStr)
	r, err = os.Open(pathStr)
	if err != nil {
		name = fmt.Sprintf("%s.txt", name) // force to .txt
		r = errReadCloser(err)
	}

	return name, r, nil
}
