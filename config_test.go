package goload

import (
	"testing"
)

func TestConfig(t *testing.T) {
	task := NewTask("/test", GET)
	if task.Path != "/test" && task.Method != GET {
		t.Errorf("not match path")
	}

	type TestPost struct {
		name string
	}

	err := make(chan error)
	go func() {
		defer close(err)
		err <- task.SetQuery("test_int", 10)
		err <- task.SetQuery("test_float", 10.0)
		err <- task.SetQuery("test_string", "test")
    err <- task.SetQuery("test_bool", true)

		err <- task.SetPostData("test_int", 10)
		err <- task.SetPostData("test_float", 10.0)
		err <- task.SetPostData("test_string", "test")
    err <- task.SetPostData("test_bool", false)

	}()

	for v := range err {
		if v != nil {
			t.Errorf("failed set", v)
		}
	}

}
