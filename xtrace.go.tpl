package {{if .UniqueString}}xtracego_{{.UniqueString}}{{else}}main{{end}}

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

func getTimestamp() string {
	return time.Now().In(time.UTC).Format(time.RFC3339)
}

func getGoroutineId() string {
	return string(bytes.Split(debug.Stack(), []byte(" "))[1])
}

func getFuncName() string {
	pc, _, _, _ := runtime.Caller(2) // > PrintlnCall > getFuncName
	return runtime.FuncForPC(pc).Name()
}

func getPrefix(showTimestamp, showGoroutine bool) string {
	prefix := ""
	if showTimestamp {
		prefix += fmt.Sprintf("%20s ", getTimestamp())
	}
	if showGoroutine {
		prefix += fmt.Sprintf("[%2s] ", getGoroutineId())
	}
	return prefix + ":"
}

func PrintlnStatement_{{.UniqueString}}(line string, showTimestamp, showGoroutine bool) {
	_, _ = fmt.Fprintln(os.Stderr, getPrefix(showTimestamp, showGoroutine)+line)
}

func PrintlnVariable_{{.UniqueString}}(varName string, varValue any, showTimestamp, showGoroutine bool) {
	_, _ = fmt.Fprintf(os.Stderr, getPrefix(showTimestamp, showGoroutine)+"[VAR] "+varName+"=%#v\n", varValue)
}

func PrintlnCall_{{.UniqueString}}(showTimestamp, showGoroutine bool) {
	_, _ = fmt.Fprintln(os.Stderr, getPrefix(showTimestamp, showGoroutine)+"[CALL] ("+getFuncName()+")")
}

func PrintlnReturn_{{.UniqueString}}(showTimestamp, showGoroutine bool) {
	_, _ = fmt.Fprintln(os.Stderr, getPrefix(showTimestamp, showGoroutine)+"[RETURN] ("+getFuncName()+")")
}
