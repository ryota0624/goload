package goload

import (
  "fmt"
)

func StatSummary(statChan StatChannel, returnChan ReturnChan) {
	stats := []*Stat{}
	func() {
		for {
			select {
			case stat := <-statChan:
				stats = append(stats, stat)
			case <-returnChan:
        return
			}
		}
	}()
  
  fmt.Println(AverageResponseTime(stats))
}

func AverageResponseTime(stats []*Stat) int {
  ave := 0
  for _, stat := range stats {
    ave += stat.ResponseTime
  }

  return ave / len(stats)
}