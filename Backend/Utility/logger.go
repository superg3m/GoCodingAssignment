package Utility

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
)

var (
	COLOR_RED     = color.New(color.FgRed).SprintFunc()
	COLOR_BLUE    = color.New(color.FgBlue).SprintFunc()
	COLOR_GREEN   = color.New(color.FgGreen).SprintFunc()
	COLOR_YELLOW  = color.New(color.FgYellow).SprintFunc()
	COLOR_MAGENTA = color.New(color.FgMagenta).SprintFunc()
	COLOR_CYAN    = color.New(color.FgCyan).SprintFunc()

	COLOR_BG_RED = color.New(color.FgWhite, color.BgRed).SprintFunc()
)

func printCallStack() {
	pc := make([]uintptr, 32)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var allFrames []runtime.Frame
	for {
		frame, more := frames.Next()
		allFrames = append(allFrames, frame)
		if !more {
			break
		}
	}

	LogPrint(COLOR_CYAN("Call Stack:"))
	for i := 0; i < len(allFrames)-1; i++ {
		frame := allFrames[i]
		relativePath, pathErr := filepath.Rel(".", frame.File)
		if pathErr != nil {
			relativePath = frame.File
		}

		fmt.Printf("%s:%d\n", COLOR_MAGENTA(relativePath), frame.Line)
	}
}

func log(msg string) {
	fmt.Println(msg)
}

func LoggerInit() {
	// stdout OR filePtr
}

// TODO (Jovanni): need to fix this!
func LogOnError(err error, format string, args ...any) {
	if err != nil {
		LogFatal(fmt.Sprintf("Error: %s\n", err))
		panic(err)
	}
}

func LogPrint(msg string, args ...any) {
	log(fmt.Sprintf(msg, args...))
}

func LogSuccess(msg string, args ...any) {
	log(COLOR_GREEN(fmt.Sprintf(msg, args...)))
}

func LogDebug(msg string, args ...any) {
	log(COLOR_BLUE(fmt.Sprintf(msg, args...)))
}

func LogWarn(msg string, args ...any) {
	log(COLOR_MAGENTA(fmt.Sprintf(msg, args...)))
}

func LogError(msg string, args ...any) {
	log(COLOR_RED(fmt.Sprintf(msg, args...)))
}

func LogFatal(msg string, args ...any) {
	log(COLOR_BG_RED(fmt.Sprintf(msg, args...)))
}
