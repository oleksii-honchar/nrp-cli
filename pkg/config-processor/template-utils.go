package configProcessor

import (
	c "beaver/blablo/color"
	"os"
	"text/template"
)

// var f = fmt.Sprintf
var confTemplate *template.Template
var defaultConfTemplate *template.Template

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

	logger.Info(f("Template loaded: %s", c.WithGreen(filePath)))

	return tmpl, nil
}

func loadConfTemplate(filePath string) (bool, error) {
	var err error
	confTemplate, err = loadTemplate(filePath)
	return err == nil, err
}

func loadDefaultConfTemplate(filePath string) (bool, error) {
	var err error
	defaultConfTemplate, err = loadTemplate(filePath)
	return err == nil, err
}
