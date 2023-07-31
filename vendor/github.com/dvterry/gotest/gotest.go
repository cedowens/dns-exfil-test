package gotest

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

const (
	resetColor = "\033[0m"
	red        = "\033[31m"
	green      = "\033[32m"
	yellow     = "\033[33m"
	blue       = "\033[34m"
	purple     = "\033[35m"
	cyan       = "\033[36m"
	gray       = "\033[37m"
	white      = "\033[97m"
)

func ErrIsNil(t testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(red+"%s:%d unexpected error: %s\n"+resetColor, filepath.Base(file), line, err)
		t.FailNow()
	}
}

func IsEqual(t testing.TB, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(red+"%s:%d unequal values: got: %#v want: %#v\n"+resetColor, filepath.Base(file), line, got, want)
		t.FailNow()
	}
}

func IsTrue(t testing.TB, condition bool, msg string) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(red+"%s:%d "+msg+"\n"+resetColor, filepath.Base(file), line)
		t.FailNow()
	}
}
