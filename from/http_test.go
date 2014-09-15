package from

import "testing"
import "github.com/nowk/assert"

func TestHTTPURLParseError(t *testing.T) {
	badURL := "thisisabadurl"
	name, r, err := HTTP(badURL)

	assert.Equal(t, "thisisabadurl.txt", name)
	assert.Nil(t, err)

	b := make([]byte, 32*1024)
	n, _ := r.Read(b)
	assert.Equal(t, "parse thisisabadurl: invalid URI for request", string(b[:n]))
}

func TestHTTPClientError(t *testing.T) {
	name, r, err := HTTP("http://unreachable")

	assert.Equal(t, "unreachable.txt", name)
	assert.Equal(t,
		"Get http://unreachable: dial tcp: lookup unreachable: no such host",
		err.Error())

	b := make([]byte, 32*1024)
	n, _ := r.Read(b)
	assert.Equal(t, "GET http://unreachable: http client error", string(b[:n]))
}

func TestHTTPAppendsExtFromContentType(t *testing.T) {
	t.Skip("")
}
