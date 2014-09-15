package from

import "io"
import "io/ioutil"
import "strings"

func errReadCloser(err error) (r io.ReadCloser) {
	if err != nil {
		r = ioutil.NopCloser(strings.NewReader(err.Error()))
	}

	return
}
