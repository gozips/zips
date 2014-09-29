package zips

import "bytes"
import "fmt"
import "io/ioutil"
import "net/http"
import "net/http/httptest"
import "testing"
import "github.com/nowk/assert"
import "github.com/gozips/source"

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
	zip := NewZip(source.HTTP)
	zip.Add(url1)
	zip.Add(url2, url3)
	n, err := zip.WriteTo(out)

	assert.Nil(t, err)
	assert.Equal(t, int64(38), n)
	verifyZip(t, out.Bytes(), []tEntries{
		{"index.html", "Hello World!"},
		{"posts", "Post Body"},
		{"data.json", `{"data": ["one"]}`},
	})
}

func TestZipFromFSSources(t *testing.T) {
	out := new(bytes.Buffer)
	zip := NewZip(source.FS)
	zip.Add("sample/file1.txt")
	zip.Add("sample/file2.txt")
	zip.Add("sample/file3.txt")
	n, err := zip.WriteTo(out)

	assert.Nil(t, err)
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
			r := bytes.NewReader([]byte("Good!"))
			return srcPath, ioutil.NopCloser(r)
		}

		return srcPath, fmt.Errorf("uh-oh")
	}

	out := new(bytes.Buffer)
	zip := NewZip(sourceFn)
	zip.Add("good", "error", "andgoodagain")

	_, err := zip.WriteTo(out)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("1 error(s):\n\n%s", "* uh-oh"), err.Error())

	ze := err.(ZipError)
	assert.Equal(t, 1, len(err.(ZipError)))
	assert.Equal(t, "uh-oh", ze[0].Error())
	verifyZip(t, out.Bytes(), []tEntries{
		{"good", "Good!"},
		{"andgoodagain", "Good!"},
	})
}
