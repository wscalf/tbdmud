package text

import (
	"fmt"
	"io"
)

type PrintfJob struct {
	template string
	params   []interface{}
}

func NewPrintfJob(template string, params ...interface{}) PrintfJob {
	return PrintfJob{template: template, params: params}
}

func (f PrintfJob) Run(w io.Writer) error {
	_, err := fmt.Fprintf(w, f.template, f.params...)
	return err
}
