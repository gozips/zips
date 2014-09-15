package from

import "net/url"
import "testing"
import "github.com/nowk/assert"

func TestFilename(t *testing.T) {
	for _, f := range []struct {
		Path, Basename string
	}{
		{"/path/file.txt", "file.txt"},
		{"./file.txt", "file.txt"},
		{"/file.txt", "file.txt"},
		{"file", "file"},
	} {
		name := filename(f.Path)
		assert.Equal(t, f.Basename, name)
	}

	for _, u := range []struct {
		URL      *url.URL
		Basename string
	}{
		{&url.URL{Host: "example.com", Path: ""}, "example.com"},
		{&url.URL{Host: "example.com", Path: "/"}, "example.com"},
		{&url.URL{Path: "/foo"}, "foo"},
		{&url.URL{Path: "/foo/"}, "foo"},
		{&url.URL{Path: "/foo.html"}, "foo.html"},
		{&url.URL{Path: "/foo.json"}, "foo.json"},
	} {
		name := filename(u.URL)
		assert.Equal(t, u.Basename, name)
	}
}
