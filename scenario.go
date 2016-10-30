package goload

import (
	"net/url"
	_ "reflect"
	"strconv"
	"strings"
)

type RequestMethod int

const (
	GET RequestMethod = iota + 1
	POST
	PUT
	DELETE
)

func StrToRequestMethod(str string) (RequestMethod, error) {
	upperStr := strings.ToUpper(str)
	switch upperStr {
		case "GET":
			return GET, nil
		case "POST":
			return POST, nil
		case "PUT":
			return PUT, nil
		case "DELETE":
			return DELETE, nil
		default:
			return GET, InvalidValue("")
	}
}

type URLData map[string]string

type InvalidValue string

func (e InvalidValue) Error() string {
	return "InvalidValue"
}

func (d URLData) SetData(name string, value interface{}) error {
	str, err := func() (string, error) {
		switch r := value.(type) {
		case int:
			return strconv.Itoa(int(r)), nil
		case string:
			return string(r), nil
		case float64:
			return strconv.FormatFloat(r, 'f', 4, 64), nil
		case bool:
			return strconv.FormatBool(r), nil
		default:
			return "", InvalidValue("error")
		}
	}()

	d[name] = str

	return err
}

func emptyURLData() URLData {
	return URLData(map[string]string{})
}

type Task struct {
	Method   RequestMethod
	Path     string
	bodyType     string
	postData URLData
	query    URLData
}

type Scenario struct {
	tasks []Task
}

func NewTask(path string, method RequestMethod) Task {
	return Task{
		Path:   path,
		Method: method,
	}
}

func (s *Task) SetBody(bodyType string) {
	s.bodyType = bodyType
}

func (s *Task) SetQuery(name string, value interface{}) error {
	if len(s.query) == 0 {
		s.query = emptyURLData()
	}
	return s.query.SetData(name, value)
}

func (s *Task) SetPostData(name string, value interface{}) error {
	if len(s.postData) == 0 {
		s.postData = emptyURLData()
	}
	return s.postData.SetData(name, value)
}

func (s *Task) URL(hostname string) *url.URL {
	var (
		err error
		urlBase *url.URL
	)

	if urlBase, err = url.Parse(hostname + s.Path); err != nil {
		panic(err)
	}

	if len(s.query) != 0 {
		query := urlBase.Query()
		for k, v := range s.query {
			query.Set(k, v)
		}
		urlBase.RawQuery = query.Encode()
	}
	return urlBase
}

func (s *Task) PostData() url.Values {
	data := url.Values{}
	if len(s.postData) != 0 {
		for k, v := range s.postData {
			data.Set(k, v)
		}
	}
	return data
}
