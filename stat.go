package goload

import (
  "time"
  // "strconv"
  "text/template"
  "io"
  _"fmt"
)

type Stat struct {
  Path string
  ResponseTime int
  StatusCode int
  StartTime time.Time
  EndTime time.Time
  Err error
}

var statTemplate = template.New("statTemplate")
const statTamplateStr = "{{.Path}},{{.ResponseTime}},{{.StatusCode}},{{.StartTime}},{{.EndTime}}\n"
func init() {
  var err error
  if statTemplate, err = statTemplate.Parse(statTamplateStr); err != nil {
    panic(err)
  }
}

func (s Stat) Write(io io.Writer) error {
  return statTemplate.Execute(io, s)
}

func NewStat(path string) Stat {
  return Stat {
    Path: path,
    StartTime: time.Now(),
  }
}

func (s *Stat) SetError(err error) {
  s.Err = err
}

func (s *Stat) FinishMesure(statusCode int) {
  s.EndTime = time.Now()
  s.StatusCode = statusCode
  s.ResponseTime = int(s.EndTime.UnixNano() - s.StartTime.UnixNano())
}