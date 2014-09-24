package zips_test

import "bytes"
import "errors"
import "fmt"
import "io/ioutil"
import "net/http"
import "net/http/httptest"
import "strings"
import "testing"
import "github.com/nowk/assert"
import . "github.com/nowk/go-zips"
import "github.com/nowk/go-zips/from"

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
	zip := NewZip(url1, url2, url3)
	n, ok := zip.Write(out, from.HTTP)

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
	zip := NewZip("sample/file1.txt", "sample/file2.txt", "sample/file3.txt")
	n, ok := zip.Write(out, from.FS)

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
	out := new(bytes.Buffer)
	zip := NewZip("good", "error", "andgoodagain")
	_, ok := zip.Write(out, func(srcPath string) (string, interface{}) {
		if "good" == srcPath || "andgoodagain" == srcPath {
			return srcPath, ioutil.NopCloser(strings.NewReader("Good!"))
		}

		return srcPath, errors.New("uh-oh")
	})

	assert.False(t, ok)
	assert.Equal(t, 1, len(zip.Errors()))
	assert.Equal(t, "uh-oh", zip.Errors()[0].Error())
	verifyZip(t, out.Bytes(), []tEntries{
		{"good", "Good!"},
		{"andgoodagain", "Good!"},
	})
}

func TestAdd(t *testing.T) {
	zip := NewZip()
	zip.Add("source1.txt")
	zip.Add("source2.txt")
	zip.Add("source3.txt")
	assert.Equal(t, 3, len(zip.Sources))

	zip = NewZip("source1.txt", "source2.txt")
	zip.Add("source3.txt")
	zip.Add("source4.txt")
	zip.Add("source5.txt")
	assert.Equal(t, 5, len(zip.Sources))
}
