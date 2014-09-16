package from

import "io"
import "regexp"
import "testing"
import "github.com/nowk/assert"

func TestHTTPURLParseError(t *testing.T) {
	badURL := "thisisabadurl"
	name, v := HTTP(badURL)
	r := v.(io.ReadCloser)
	defer r.Close()

	assert.Equal(t, "thisisabadurl.txt", name)

	b := make([]byte, 32*1024)
	n, _ := r.Read(b)
	assert.Equal(t, "parse thisisabadurl: invalid URI for request", string(b[:n]))
}

func TestHTTPClientError(t *testing.T) {
	name, v := HTTP("http://unreachable")
	r := v.(io.ReadCloser)
	defer r.Close()

	assert.Equal(t, "unreachable.txt", name)

	reg := regexp.MustCompile(`Get http:\/\/unreachable:( dial tcp:)? lookup unreachable: no such host`)
	b := make([]byte, 32*1024)
	n, _ := r.Read(b)
	if str := string(b[:n]); !reg.MatchString(str) {
		t.Errorf("Expected %s, got %s", reg.String(), str)
	}
}

// func TestHTTPAppendsExtFromContentType(t *testing.T) {
// 	t.Skip("")
// }
