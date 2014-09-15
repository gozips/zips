package from

import "testing"
import "github.com/nowk/assert"

func TestFSPathError(t *testing.T) {
	name, r, err := FS("path/to/doesnotexist.txt")
	defer r.Close()

	assert.Equal(t, "doesnotexist.txt.txt", name)
	assert.Nil(t, err)

	b := make([]byte, 32*1024)
	n, _ := r.Read(b)
	assert.Equal(t, "open path/to/doesnotexist.txt: no such file or directory", string(b[:n]))
}
