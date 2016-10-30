package goload

import (
  "fmt"
  "bufio"
  "log"
  "os"
)
type OUTPUT interface {
  Close()
  Write(p []byte) (int, error)
}
type CLIOUTPUT struct {}
func (o CLIOUTPUT) Write(p []byte) (int, error) {
  fmt.Print(string(p))
  return len(p), nil
}
func (o CLIOUTPUT) Close() {

}

func newFile(fn string) *os.File {
    fp, err := os.Create(fn)
    if err != nil {
        log.Fatal(err)
    }
    return fp
}


type FileOUTPUT struct {
  Path string
  File *os.File
  Writer *bufio.Writer
}

func NewFileOutPUT(filename string) FileOUTPUT {
  file := newFile(filename)
  return FileOUTPUT {
    Path: filename,
    File: file,
    Writer: bufio.NewWriter(file),
  }
}

func (f FileOUTPUT) Write(p []byte) (int, error) {
  return f.Writer.Write(p)
}

func (o FileOUTPUT) Close() {
  var err error
  err = o.Writer.Flush()
  err = o.File.Close()
  if err != nil {
    panic(err)
  }
}