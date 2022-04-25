package logs

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MyLogFormatter struct{}

func (m *MyLogFormatter) Format(entry *log.Entry) ([]byte, error) {
	stamp := entry.Time.Format("2006-01-02 15:04:05.000")
	levelStr := strings.ToUpper(entry.Level.String())
	pid := os.Getpid()
	fileName := entry.Data["__file__"]
	lineNo := entry.Data["__line__"]
	b := entry.Buffer
	if b == nil {
		b = &bytes.Buffer{}
	}
	fmt.Fprintf(b, "%s [%s][%d] %v %v: %v", stamp, levelStr, pid, fileName, lineNo, entry.Message)
	appendKVsAndNewLine(b, entry)
	return b.Bytes(), nil
}

func appendKVsAndNewLine(b *bytes.Buffer, entry *log.Entry) {
	// Sort the keys for consistent output.
	var keys []string = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if key == "__file__" || key == "__line__" {
			continue
		}
		var value interface{} = entry.Data[key]
		var stringifiedValue string
		if err, ok := value.(error); ok {
			stringifiedValue = err.Error()
		} else if stringer, ok := value.(fmt.Stringer); ok {
			// Trust the value's String() method.
			stringifiedValue = stringer.String()
		} else {
			// No string method, use %#v to get a more thorough dump.
			fmt.Fprintf(b, " %v=%#v", key, value)
			continue
		}
		b.WriteByte(' ')
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteString(stringifiedValue)
	}
	b.WriteByte('\n')
}

type ContextHook struct {
}

func (hook ContextHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook ContextHook) Fire(entry *log.Entry) error {
	pcs := make([]uintptr, 4)
	if numEntries := runtime.Callers(6, pcs); numEntries > 0 {
		frames := runtime.CallersFrames(pcs)
		for {
			frame, more := frames.Next()
			if !shouldSkipFrame(frame) {
				entry.Data["__file__"] = path.Base(frame.File)
				entry.Data["__line__"] = frame.Line
				break
			}
			if !more {
				break
			}
		}
	}
	return nil
}

func shouldSkipFrame(frame runtime.Frame) bool {
	return strings.LastIndex(frame.File, "exported.go") > 0 ||
		strings.LastIndex(frame.File, "logger.go") > 0 ||
		strings.LastIndex(frame.File, "entry.go") > 0
}

func LoggerInit() {
	log.SetOutput(os.Stdout)
	//log.SetLevel(log.DebugLevel)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(new(MyLogFormatter))
	log.AddHook(new(ContextHook))
}