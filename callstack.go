package callstack

import (
	"encoding/json"
	"runtime"
	"strings"
)

const callStackDepth = 50 // most relevant context seem to appear near the top of the stack

// Frame -
type Frame struct {
	File string `json:"file"`
	Line int    `json:"line"`
	Func string `json:"function"`
}

// String - `stringer` interface
func (thisRef Frame) String() string {
	data, _ := json.Marshal(thisRef)
	return string(data)
}

// Get - Gets the call stack with no frames skipped
func GetCurrent() (packageName string, stackFrames []Frame) {
	return GetCurrentSkipFrames(0)
}

// GetSkipFrames - Gets the call stack with N skipped farames
func GetCurrentSkipFrames(skip int) (packageName string, stackFrames []Frame) {
	var callStack [callStackDepth]uintptr
	callStackSize := runtime.Callers(skip, callStack[:])
	return getFromExistingCallers(callStack[:callStackSize])
}

func getFromExistingCallers(callStack []uintptr) (packageName string, stackFrames []Frame) {
	packageName = ""
	stackFrames = []Frame{}

	callStackFrames := runtime.CallersFrames(callStack)
	for {
		frame, ok := callStackFrames.Next()
		if !ok {
			break
		}

		pkg, fn := splitPackageFuncName(frame.Function)
		if packageName == "" && pkg != "runtime" { // pickup first package
			packageName = pkg
		}

		if frameFilter(pkg, fn, frame.File, frame.Line) {
			stackFrames = stackFrames[:0]
			continue
		}

		stackFrames = append(stackFrames, Frame{
			File: frame.File,
			Line: frame.Line,
			Func: fn,
		})
	}

	return
}

func splitPackageFuncName(funcName string) (string, string) {
	var packageName string
	if ind := strings.LastIndex(funcName, "/"); ind > 0 {
		packageName += funcName[:ind+1]
		funcName = funcName[ind+1:]
	}
	if ind := strings.Index(funcName, "."); ind > 0 {
		packageName += funcName[:ind]
		funcName = funcName[ind+1:]
	}
	return packageName, funcName
}

func frameFilter(packageName, funcName string, file string, line int) bool {
	return packageName == "runtime" && funcName == "panic"
}
