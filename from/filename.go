package from

import "path/filepath"
import "net/url"

// filename parses the base name from a url.URL or string (for fs)
func filename(o interface{}) string {
	var name string

	switch v := o.(type) {
	case string:
		return filepath.Base(v)
	case *url.URL:
		name = filepath.Base(v.Path)
		if "." == name || "/" == name {
			name = v.Host
		}
	}

	return name
}
