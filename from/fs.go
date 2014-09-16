package from

import "fmt"
import "os"

// FS returns the file name, and either a ReadCloser or error as interface{}
func FS(pathStr string) (string, interface{}) {
	name := filename(pathStr)
	r, err := os.Open(pathStr)
	if err != nil {
		name = fmt.Sprintf("%s.txt", name) // force to .txt
		return name, errReadCloser(err)
	}

	return name, r
}
