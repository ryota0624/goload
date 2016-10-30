package goload

import (
	"testing"
	// "time"
)

func TestExecuter(t *testing.T) {
	task := NewTask("/", GET)
	tasks := []Task{task}
	jobFinChan := make(chan bool)
	ch := make(StatChannel)
	scenario := Scenario{
		tasks: tasks,
	}
	go ScenarioExecuter("http://localhost:8080", scenario, ch, jobFinChan)

	results := []*Stat{}
	func() {
		for {
			select {
			case <-jobFinChan:
				return
			case stat := <-ch:
				results = append(results, stat)
			}
		}
	}()

	if len(results) != 1 {
		t.Errorf("not match resultStatlength")
	}
}

func BenchmarkExecuter(t *testing.B) {
	for i := 0; i < t.N; i++ {
		task := NewTask("/", GET)
		tasks := []Task{task}
		jobFinChan := make(chan bool)
		ch := make(StatChannel)
		scenario := Scenario{
			tasks: tasks,
		}

		go ScenarioExecuter("http://localhost:8080", scenario, ch, jobFinChan)

		results := []*Stat{}
		func() {
			for {
				select {
				case <-jobFinChan:
					return
				case stat := <-ch:
					results = append(results, stat)
				}
			}
		}()
		if len(results) != 1 {
			t.Errorf("not match resultStatlength")
		}
	}
}

func TestRun(t *testing.T) {
	tasks := []Task{}
	ch := make(StatChannel)
	jobFinishChan := make(chan bool)

	for i := 0; i < 10; i++ {
		tasks = append(tasks, NewTask("/", GET))
	}
	scenario := Scenario{
		tasks: tasks,
	}
	// for i := 0; i < 10; i++ {
	//   task := NewTask("/", POST)
	//   task.SetBody("application/x-www-form-urlencoded")
	//   task.SetPostData("text", "golang")
	// 	tasks = append(tasks, task)
	//

	rate := 10
	duration := 3
	conf := Conf{
		Scenario: scenario,
		Target:   "http://localhost:8080",
		Duration: duration,
		Rate:     rate,
	}
	go Run(conf, ch, jobFinishChan)

	statsArr := []*Stat{}

	func() {
		for {
			select {
			case stat := <-ch:
				statsArr = append(statsArr, stat)
			case <-jobFinishChan:
				return
			}
		}
	}()

	if len(statsArr) != rate*duration*len(tasks) {
		t.Errorf("stats length no mathc%v", len(statsArr))
	}
}
