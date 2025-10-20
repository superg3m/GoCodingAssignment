package Utility

import (
	"fmt"
	"os"
)

func AssertOnError(err error) {
	if err != nil {
		LogFatal(fmt.Sprint("[Developer Error]: ", err))
		printCallStack()
		os.Exit(-1)
	}
}

func AssertOnErrorMsg(err error, format string, args ...any) {
	if err != nil {
		combined := fmt.Sprint("[Developer Error]: ", format)
		LogFatal(fmt.Sprintf(combined, args...))
		printCallStack()
		os.Exit(-1)
	}
}

func Assert(condition bool) {
	if !condition {
		LogFatal("[Assert Triggered]")
		printCallStack()
		os.Exit(-1)
	}
}

func AssertMsg(condition bool, format string, args ...any) {
	if !condition {
		msg := fmt.Sprintf(format, args...)
		msg2 := format
		if len(args) != 0 {
			LogFatal(fmt.Sprint("[Assert Triggered]: ", msg))
		} else {
			LogFatal(fmt.Sprint("[Assert Triggered]: ", msg2))
		}
		printCallStack()
		os.Exit(-1)
	}
}
