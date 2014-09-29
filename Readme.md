# go-zips

[![Build Status](https://travis-ci.org/gozips/zips.svg?branch=master)](https://travis-ci.org/gozips/zips)
[![GoDoc](https://godoc.org/github.com/gozips/zips?status.svg)](http://godoc.org/github.com/gozips/zips)

An API to always return a zip archive

## Example

    import "log"
    import "os"
    import "github.com/gozips/zips"
    import "github.com/gozips/sources"

    func main() {
      out := os.Create("out.zip")
      zip := zips.NewZip(sources.FS)
      zip.Add("file1.txt")
      zip.Add("file2.txt")
      zip.Add("file3.txt")
      n, ok := zip.WriteTo(out)

      log.Print("bytes written ", n)
      log.Print("archived without error ", ok)
    }

## License

MIT
