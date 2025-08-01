package stacktrace

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const maxStacktraceDepth = 32

type Frame uintptr

func (f Frame) pc() uintptr {
	return uintptr(f) - 1
}

func (f Frame) String() string {
	funcForPC := runtime.FuncForPC(f.pc())
	if funcForPC == nil {
		return ""
	}

	_, name := filepath.Split(funcForPC.Name())

	file, line := funcForPC.FileLine(f.pc())
	dir, file := filepath.Split(file)
	pathParts := strings.SplitAfter(dir, string(os.PathSeparator))

	var lastDir string

	if len(pathParts) > 0 {
		lastDir = pathParts[len(pathParts)-2]
	}

	file = lastDir + file

	return fmt.Sprintf("%s\n%s:%d", name, file, line)
}

type StackTrace []Frame

func (s StackTrace) String() string {
	sb := strings.Builder{}

	for i, frame := range s {
		sb.WriteString(frame.String())
		if i < len(s)-1 {
			sb.WriteRune('\n')
		}
	}

	return sb.String()
}

// Trace возвращает стектрейс глубиной не более maxStacktraceDepth.
// Возвращаемый стектрейс начинается с того места, где была вызвана Trace.
func Trace() StackTrace {
	programCounters := make([]uintptr, maxStacktraceDepth)
	callers := runtime.Callers(2, programCounters)
	programCounters = programCounters[:callers]

	stackTrace := StackTrace{}

	for _, pc := range programCounters {
		stackTrace = append(stackTrace, Frame(pc))
	}

	return stackTrace
}
