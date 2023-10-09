package nginxConfigProcessor

import (
	c "github.com/oleksii-honchar/coteco"

	"os"
	"text/template"
)

var confTemplate *template.Template
var defaultConfTemplate *template.Template
var acmeChallengeTemplate *template.Template
var corsServersTemplate *template.Template

func loadTemplate(filePath string) (*template.Template, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("Error reading template file", "err", err)
		return nil, err
	}

	tmpl, err := template.New(filePath).Parse(string(content))
	if err != nil {
		logger.Error("Error parsing template", "err", err)
		return nil, err
	}

	logger.Debug(f("Template loaded: %s", c.WithGreen(filePath)))

	return tmpl, nil
}

func loadTemplateToVar(filePath string, destVar **template.Template) (bool, error) {
	var err error
	*destVar, err = loadTemplate(filePath)
	return err == nil, err
}
