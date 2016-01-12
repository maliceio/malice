package errors

import (
	"fmt"
	"runtime"

	"github.com/Sirupsen/logrus"
)

// CheckError check error and return true if error is nil and return false if error is not nil.
func CheckError(err error) bool {
	return CheckErrorWithMessage(err, "")
}

// CheckErrorNoStack check error and return true if error is nil and return false if error is not nil.
func CheckErrorNoStack(err error) bool {
	return CheckErrorNoStackWithMessage(err, "")
}

// CheckErrorWithMessage check error with message and log messages with stack. And then return true if error is nil and return false if error is not nil.
func CheckErrorWithMessage(err error, msg string, args ...interface{}) bool {
	if err != nil {
		var stack [4096]byte
		runtime.Stack(stack[:], false)
		//		logrus.Errorf(msg, args)
		if len(args) == 0 {
			logrus.Error(msg + fmt.Sprintf("%q\n%s\n", err, stack[:]))
		} else {
			logrus.Error(fmt.Sprintf(msg, args...) + fmt.Sprintf("%q\n%s\n", err, stack[:]))
		}
		//		logrus.Printf(msg+"\n%q\n%s\n",args, err, stack[:])
		return false
	}
	return true
}

// CheckErrorNoStackWithMessage check error with message and return true if error is nil and return false if error is not nil.
func CheckErrorNoStackWithMessage(err error, msg string, args ...interface{}) bool {
	if err != nil {
		if len(args) == 0 {
			logrus.Error(msg + fmt.Sprintf("%q\n", err))
		} else {
			logrus.Error(fmt.Sprintf(msg, args...) + fmt.Sprintf("%q\n", err))
		}
		//		logrus.Printf(msg+"\n%q\n%s\n",args, err, stack[:])
		return false
	}
	return true
}
