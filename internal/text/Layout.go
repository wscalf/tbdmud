package text

import (
	"strings"
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

func (l *Layout) Prepare(obj Formattable) FormatJob {
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

func (t LayoutJob) Run() (string, error) {
	var builder strings.Builder
	err := t.template.Execute(&builder, t.properties)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}
