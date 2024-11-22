package safe

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func Go(fn func() error) {
	go func() {
		defer Recover()

		if err := fn(); err != nil {
			logrus.Errorf("error: %v", err)
		}
	}()
}

func Recover() {
	e := recover()
	if e == nil {
		return
	}

	debug.PrintStack()

	if err, ok := e.(error); ok {
		logrus.Errorf("panic error: %+v", err)
	} else {
		logrus.Errorf("panic error: %+v", e)
	}
}
