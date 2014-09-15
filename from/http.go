package from

import "fmt"
import "io"
import "net/http"
import "net/url"

func HTTP(urlStr string) (string, io.ReadCloser, error) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Sprintf("%s.txt", urlStr), errReadCloser(err), nil
	}

	name := filename(u)
	req := &http.Request{Method: "GET", URL: u}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Sprintf("%s.txt", name),
			errReadCloser(fmt.Errorf("%s %s: http client error", req.Method, urlStr)),
			err
	}

	return name, resp.Body, nil
}
