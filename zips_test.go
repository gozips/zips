package zips

import (
	"bytes"
	"fmt"
	"github.com/gozips/source"
	"github.com/gozips/sources"
	gozipst "github.com/gozips/testing"
	"github.com/nowk/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestErrorIsTrueNil(t *testing.T) {
	out := new(bytes.Buffer)
	zip := NewZip(sources.HTTP)
	_, err := zip.WriteTo(out)

	if err != nil {
		t.Error("expected true nil")
	}
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
	zip.Add(url2, url3)
	n, err := zip.WriteTo(out)

	assert.Nil(t, err)
	assert.Equal(t, int64(38), n)
	assert.Equal(t, int64(38), zip.BytesIn)
	assert.Equal(t, int64(56), zip.BytesOut)
	gozipst.VerifyZip(t, out.Bytes(), []gozipst.Entries{
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
	n, err := zip.WriteTo(out)

	assert.Nil(t, err)
	assert.Equal(t, int64(11), n)
	assert.Equal(t, int64(11), zip.BytesIn)
	assert.Equal(t, int64(29), zip.BytesOut)
	gozipst.VerifyZip(t, out.Bytes(), []gozipst.Entries{
		{"file1.txt", "One"},
		{"file2.txt", "Two"},
		{"file3.txt", "Three"},
	})
}

func tSourceFunc(c io.ReadCloser) source.Func {
	return func(srcPath string) (string, io.ReadCloser, error) {
		if "good" == srcPath || "andgoodagain" == srcPath {
			r := bytes.NewReader([]byte("Good!"))
			return srcPath, ioutil.NopCloser(r), nil
		}
		return srcPath, c, fmt.Errorf("uh-oh")
	}
}

func TestEntrySkippedIfReadCloserIsNilOnError(t *testing.T) {
	sourceFn := tSourceFunc(nil)

	out := new(bytes.Buffer)
	zip := NewZip(sourceFn)
	zip.Add("good", "error", "andgoodagain")

	_, err := zip.WriteTo(out)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("1 error(s):\n\n%s", "* uh-oh"), err.Error())

	ze := err.(Error)
	assert.Equal(t, 1, len(err.(Error)))
	assert.Equal(t, "uh-oh", ze[0].Error())
	gozipst.VerifyZip(t, out.Bytes(), []gozipst.Entries{
		{"good", "Good!"},
		{"andgoodagain", "Good!"},
	})
}

func TestEntryCreatedIfReadCloserIsNotNilOnError(t *testing.T) {
	c := ioutil.NopCloser(bytes.NewReader([]byte("uh-oh")))
	sourceFn := tSourceFunc(c)

	out := new(bytes.Buffer)
	zip := NewZip(sourceFn)
	zip.Add("good", "error", "andgoodagain")

	_, err := zip.WriteTo(out)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("1 error(s):\n\n%s", "* uh-oh"), err.Error())

	ze := err.(Error)
	assert.Equal(t, 1, len(err.(Error)))
	assert.Equal(t, "uh-oh", ze[0].Error())
	gozipst.VerifyZip(t, out.Bytes(), []gozipst.Entries{
		{"good", "Good!"},
		{"error", "uh-oh"},
		{"andgoodagain", "Good!"},
	})
}
