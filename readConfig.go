package goload

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type TaskDto struct {
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
}

type ConfDto struct {
	Target   string `yaml:"target"`
	Duration int    `yaml:"duration"` // n秒間
	Rate     int    `yaml:"rate"`     // 1秒/n req
  Scenario []TaskDto `yaml:"scenario"`
}

func ReadConfig(path string) (Conf, error) {
	var err error
	var buf []byte
	buf, err = ioutil.ReadFile(path)
	var confDto ConfDto
	if err != nil {
		return Conf{}, err
	}
	err = yaml.Unmarshal([]byte(buf), &confDto)
	if err != nil {
		return Conf{}, err
	}
	return DtoToConf(confDto)
}

func DtoToConf(dto ConfDto) (Conf, error) {
	tasks := make([]Task, len(dto.Scenario))
	for k, taskDto := range dto.Scenario {
		var err error
		if tasks[k], err = DtoToTask(taskDto); err != nil {
			return Conf{}, nil
		}
	}

	return Conf{
		Target:   dto.Target,
		Duration: dto.Duration,
		Rate:     dto.Rate,
		Scenario: Scenario { tasks: tasks },
	}, nil
}

func DtoToTask(dto TaskDto) (Task, error) {
	var (
		err error
		method RequestMethod
	)
	if method, err = StrToRequestMethod(dto.Method); err != nil {
		return Task{}, nil
	}
	
	return Task {
		Path: dto.Path,
		Method: method,
	}, nil
}