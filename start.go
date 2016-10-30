package goload

import (
  _"fmt"
)

type Args struct {
  ConfigFileName string
  OutputFilename string
}

func Start(args Args) error {
  conf, err := ReadConfig(args.ConfigFileName)
  statChan := make(StatChannel)
  returnChan := make(ReturnChan)
  if err != nil {
    panic(err)
  }
  go Run(conf, statChan, returnChan)

  var output OUTPUT

  if len(args.OutputFilename) > 0 {
    output = NewFileOutPUT(args.OutputFilename)
  } else {
    output = CLIOUTPUT{}
  }
  // go StatSummary(statChan, returnChan)

  defer output.Close()

  return func() error {
    for {
      select {
        case stat := <- statChan:
          err := stat.Write(output)
          if err != nil {
            return err
          }
        case <- returnChan:
          return nil
      }
    }
  }()
}