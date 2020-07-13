package callstack

import (
	"encoding/json"
	"runtime"
	"strings"
)

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

// GetFrames - Gets the call stack with no frames skipped
func GetFrames() []Frame {
	return GetFramesWithSkip(0)
}

// GetFramesWithSkip - Gets the call stack with N skipped frames
func GetFramesWithSkip(skip int) []Frame {
	return GetFramesFromRawFrames(GetRawFrames(skip))
}

// GetRawFrames -
func GetRawFrames(skip int) []uintptr {
	const callStackDepth = 50 // most relevant context seem to appear near the top of the stack
	var callStackBuffer = make([]uintptr, callStackDepth)
	callStackSize := runtime.Callers(skip, callStackBuffer)
	return callStackBuffer[:callStackSize]
}

// GetFramesFromRawFrames -
func GetFramesFromRawFrames(callStack []uintptr) []Frame {
	frames := []Frame{}

	callStackFrames := runtime.CallersFrames(callStack)
	for {
		frame, ok := callStackFrames.Next()
		if !ok {
			break
		}

		pkg, fn := splitPackageFuncName(frame.Function)
		if frameFilter(pkg, fn, frame.File, frame.Line) {
			frames = frames[:0]
			continue
		}

		frames = append(frames, Frame{
			File: frame.File,
			Line: frame.Line,
			Func: fn,
		})
	}

	return frames
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
