package internal

/*
Utils getting from https://github.com/benbjohnson/testing
*/

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var prefix = os.Getenv("CONFIG_PREFIX")

var DefaultEnvs = map[string]string{
	prefix + "_REDIS_STRING":   "localhost:6379",
	prefix + "_REDIS_USERNAME": "admin",
	prefix + "_REDIS_PASSWORD": "password",
	prefix + "_API_STRING":     "localhost:8080",
	prefix + "_DB_STRING":      "postgres://admin:admin@localhost/admin",
	prefix + "_DB_PROVIDER":    "psql",
	prefix + "_JWT_SECRET":     "supersecretkey",
}

// TestAssert fails the test if the condition is false.
func TestAssert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// TestOk fails the test if an err is not nil.
func TestOk(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// TestEqual fails the test if exp is not equal to act.
func TestEqual(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func DefaultErrHandler(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func EnvSetup(envValues map[string]string) {
	if envValues == nil {
		envValues = DefaultEnvs
	}
	for k, el := range envValues {
		_ = os.Setenv(k, el)
	}
}

func EnvRemove(envValues map[string]string) error {
	if envValues == nil {
		envValues = DefaultEnvs
	}
	for k, _ := range envValues {
		if err := os.Remove(k); err != nil {
			return err
		}
	}
	return nil
}
