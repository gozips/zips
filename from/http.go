package from

import "fmt"
import "net/http"
import "net/url"

// HTTP returns the file name, and either a ReadCloser or error as interface{}
func HTTP(urlStr string) (string, interface{}) {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Sprintf("%s.txt", urlStr), errReadCloser(err)
	}

	name := filename(u)
	req := &http.Request{Method: "GET", URL: u}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Sprintf("%s.txt", name), errReadCloser(err)
	}

	return name, resp.Body
}
