package text

import (
	"io"
	"text/template"
)

type Layout struct {
	template *template.Template
}

func NewLayout(name string, body string) (*Layout, error) {
	t, err := template.New(name).Parse(body)
	if err != nil {
		return nil, err
	}

	return &Layout{
		template: t,
	}, nil
}

func (l *Layout) Prepare(obj Formattable) LayoutJob {
	properties := obj.GetProperties()

	return LayoutJob{
		template:   l.template,
		properties: properties,
	}
}

type LayoutJob struct {
	template   *template.Template
	properties map[string]interface{}
}

func (t LayoutJob) Run(w io.Writer) error {
	return t.template.Execute(w, t.properties)
}
