package from

import "io"
import "testing"
import "github.com/nowk/assert"

func TestFSPathError(t *testing.T) {
	name, v := FS("path/to/doesnotexist.txt")
	r := v.(io.ReadCloser)
	defer r.Close()

	assert.Equal(t, "doesnotexist.txt.txt", name)

	b := make([]byte, 32*1024)
	n, _ := r.Read(b)
	assert.Equal(t, "open path/to/doesnotexist.txt: no such file or directory", string(b[:n]))
}
