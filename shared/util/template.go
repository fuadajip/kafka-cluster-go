package util

import (
	"bytes"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
)

// ExecuteTemplateText will execute template text with given data
func ExecuteTemplateText(templateCode *string, templ *string, data interface{}) *string {
	t := template.Must(template.New(*templateCode).Parse(*templ))

	buf := bytes.Buffer{}
	t.Execute(&buf, data)
	return aws.String(buf.String())
}

// ExecuteTemplateFile will execute template file with given data
func ExecuteTemplateFile(filePath string, data interface{}) (string, error) {
	var tpl bytes.Buffer

	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		logger.Error(err)
		return tpl.String(), err
	}

	err = tmpl.Execute(&tpl, data)
	if err != nil {
		logger.Error(err)
		return tpl.String(), err
	}
	return tpl.String(), nil
}
