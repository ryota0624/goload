package goload

import (
	"net/http"
	"strings"
	"time"
	_"fmt"
)

func ScenarioExecuter(hostname string, scenario Scenario, ch StatChannel, returnChan chan bool) {
	var err error
	//1シナリオは 1秒
	// sleepTime := 60 / len(scenario.tasks) //n ミリsec睡眠

	for _, task := range scenario.tasks {
		stat := NewStat(task.Path)
		var response *http.Response
		func() {
			switch task.Method {
			case GET:
				response, err = http.Get(task.URL(hostname).String())
			case POST:
				response, err = http.Post(task.URL(hostname).String(), task.bodyType, strings.NewReader(task.PostData().Encode()))
			default:
				panic("error not implements http.method")
			}
			if err != nil {
				return
			}
			defer response.Body.Close()
		}()
		if err != nil {
			stat.SetError(err)
			returnChan <- false
			return
		}
		stat.FinishMesure(response.StatusCode)
		ch <- &stat
		// time.Sleep(time.Millisecond * time.Duration(sleepTime))
	}

	returnChan <- true
}

type Conf struct {
	Target   string   `yaml:"target"`
	Duration int      `yaml:"duration"` // n秒間
	Rate     int      `yaml:"rate"`     // 1秒/n req
	Scenario Scenario `yaml:"scenario"`
}

//シナリオ実行単位
type StatChannel chan *Stat
type ReturnChan chan bool

func Run(conf Conf, ch StatChannel, returnChan ReturnChan) {
	requestAmount := conf.Duration * conf.Rate
	taskFinChan := make(chan bool)
	defer close(taskFinChan)
	for i := 0; i < conf.Duration; i++ {
		for j := 0; j < conf.Rate; j++ {
			go func() {
				ScenarioExecuter(conf.Target, conf.Scenario, ch, taskFinChan)
			}()
		}
		time.Sleep(time.Second * 1) //ツールの制約としてレートは1秒ごとに設定されるため
	}
	finRequestCount := 0
	for range taskFinChan {
		if finRequestCount++; requestAmount == finRequestCount {
			returnChan <- true
			break
		}
	}
}
