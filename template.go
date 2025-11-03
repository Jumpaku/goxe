package xtracego

import (
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed xtrace.go.tpl
var xtraceGo string

var xtraceGoTemplate = template.Must(template.New("xtrace.go.tpl").Parse(xtraceGo))

type XtraceGoData struct {
	UniqueString string
}

func GetXtraceGo(uniqueString string, w io.Writer) (err error) {
	d := XtraceGoData{UniqueString: uniqueString}
	if err := xtraceGoTemplate.Execute(w, d); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
