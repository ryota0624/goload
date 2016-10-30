package goload

import (
  "testing"
)

func TestOutPut(t *testing.T) {
  stat := NewStat("/")
  stat.Write(CLIOUTPUT{})
}