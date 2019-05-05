package v1

import (
	"fmt"
	"testing"
)

func TestCheckFileExist(t *testing.T) {
	fileExist := "common_test.go"
	if err := CheckFileExist(fileExist); err != nil {
		t.Errorf("Fail test: %s", err)
	}
}

func TestLogOnError(t *testing.T) {
	LogOnError(fmt.Errorf("test error"), "error")
	LogOnError(nil, "error")
}

func TestFailOnError(t *testing.T) {
	FailOnError(fmt.Errorf("test error"), "error")
}
