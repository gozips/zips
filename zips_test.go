package zips

import "bytes"
import "errors"
import "fmt"
import "io/ioutil"
import "net/http"
import "net/http/httptest"
import "strings"
import "testing"
import "github.com/nowk/assert"
import "github.com/gozips/sources"

func h(str string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(str))
	}
}

func tServer() (ts *httptest.Server) {
	mux := http.NewServeMux()
	mux.HandleFunc("/index.html", h("Hello World!"))
	mux.HandleFunc("/posts", h("Post Body"))
	mux.HandleFunc("/api/data.json", h(`{"data": ["one"]}`))
	ts = httptest.NewServer(mux)

	return
}

func TestZipFromHTTPSources(t *testing.T) {
	ts := tServer()
	defer ts.Close()

	url1 := fmt.Sprintf("%s/index.html", ts.URL)
	url2 := fmt.Sprintf("%s/posts", ts.URL)
	url3 := fmt.Sprintf("%s/api/data.json", ts.URL)

	out := new(bytes.Buffer)
	zip := NewZip(sources.HTTP)
	zip.Add(url1)
	zip.Add(url2)
	zip.Add(url3)
	n, ok := zip.WriteTo(out)

	assert.True(t, ok)
	assert.Equal(t, 0, len(zip.Errors()))
	assert.Equal(t, int64(38), n)
	verifyZip(t, out.Bytes(), []tEntries{
		{"index.html", "Hello World!"},
		{"posts", "Post Body"},
		{"data.json", `{"data": ["one"]}`},
	})
}

func TestZipFromFSSources(t *testing.T) {
	out := new(bytes.Buffer)
	zip := NewZip(sources.FS)
	zip.Add("sample/file1.txt")
	zip.Add("sample/file2.txt")
	zip.Add("sample/file3.txt")
	n, ok := zip.WriteTo(out)

	assert.True(t, ok)
	assert.Equal(t, 0, len(zip.Errors()))
	assert.Equal(t, int64(11), n)
	verifyZip(t, out.Bytes(), []tEntries{
		{"file1.txt", "One"},
		{"file2.txt", "Two"},
		{"file3.txt", "Three"},
	})
}

func TestErrorSkipsEntry(t *testing.T) {
	var sourceFn = func(srcPath string) (string, interface{}) {
		if "good" == srcPath || "andgoodagain" == srcPath {
			return srcPath, ioutil.NopCloser(strings.NewReader("Good!"))
		}

		return srcPath, errors.New("uh-oh")
	}

	out := new(bytes.Buffer)
	zip := NewZip(sourceFn)
	zip.Add("good")
	zip.Add("error")
	zip.Add("andgoodagain")

	_, ok := zip.WriteTo(out)

	assert.False(t, ok)
	assert.Equal(t, 1, len(zip.Errors()))
	assert.Equal(t, "uh-oh", zip.Errors()[0].Error())
	verifyZip(t, out.Bytes(), []tEntries{
		{"good", "Good!"},
		{"andgoodagain", "Good!"},
	})
}
