# go-zips

[![Build Status](https://travis-ci.org/nowk/go-zips.svg?branch=master)](https://travis-ci.org/nowk/go-zips)
[![GoDoc](https://godoc.org/github.com/nowk/go-zips?status.svg)](http://godoc.org/github.com/nowk/go-zips)

An API to always return a zip archive

## go get

    go get github.com/nowk/go-zips

## Example

    import "log"
    import "os"
    import "github.com/nowk/go-zips"
    import "github.com/nowk/go-zips/from"

    func main() {
      out := os.Create("out.zip")
      zip := zips.NewZip("file1.txt", "file2.txt", "file3.txt")
      n, ok := zip.Write(from.FS, out)

      log.Print("bytes written ", n)
      log.Print("archived without error ", ok)
    }

## License

MIT
