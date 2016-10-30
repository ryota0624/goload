package goload

import (
  "testing"
)

func TestStart(t *testing.T) {
  err := Start(Args{ 
    ConfigFileName: "./conf_example.yaml", 
    OutputFilename: "./test_result.log",
   })

  if err != nil {
    t.Errorf(err.Error())
  }
}

func BenchmarkStart(b *testing.B) {
    b.StartTimer()
    Start(Args{ 
    ConfigFileName: "./conf_example.yaml", 
    OutputFilename: "./result.log",
   })
}