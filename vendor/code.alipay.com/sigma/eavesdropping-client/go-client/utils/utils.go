package utils

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"k8s.io/klog/v2"
)

// RandString generate n-length random string

func TimeSinceInMilliSeconds(start time.Time) float64 {
	return float64(time.Since(start).Nanoseconds() / time.Millisecond.Nanoseconds())
}

func TimeSinceInSeconds(start time.Time) float64 {
	return float64(time.Since(start).Nanoseconds() / time.Second.Nanoseconds())
}

// Marshal Marshal object to string of the error message if failed
func Marshal(v interface{}) string {
	bs, e := json.Marshal(v)
	if e != nil {
		return e.Error()
	}
	return string(bs)
}

// Dumps dump interface as indented json object
func Dumps(v interface{}) string {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		klog.Errorf("MarshalError(%s): %+v", err.Error(), v)
		return fmt.Sprintf("%+v", v)
	}
	return string(bs)
}

// ReTry try to call a function until success.
func ReTry(f func() error, interval time.Duration, retryCount int) (err error) {
	needContinue := true
	i := 0

	for needContinue {
		i++
		err = f()
		if err != nil && i <= retryCount {
			needContinue = true
		} else {
			needContinue = false
		}
		if needContinue {
			time.Sleep(interval)
		}
	}

	return err
}

// GetPulledImageFromEventMessage get image name from pulled event

func IgnorePanic(desc string) {
	if err := recover(); err != nil {
		klog.Error(desc, err)
		logPanic(err)
	}
}

func logPanic(r interface{}) {
	callers := getCallers(r)
	klog.Errorf("Observed a panic: %#v (%v)\n%v", r, r, callers)
}

func getCallers(r interface{}) string {
	callers := ""
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		callers = callers + fmt.Sprintf("%v:%v\n", file, line)
	}
	return callers
}
