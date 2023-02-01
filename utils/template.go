package utils

import (
	"html/template"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/sharadregoti/devops/utils/logger"
	"gopkg.in/yaml.v2"
)

// ExecuteTemplate applies a parsed template to the specified data object
func ExecuteTemplate(tmpl string, params interface{}) (string, error) {
	funcMap := map[string]interface{}{
		"toYaml": func(val interface{}) string {
			data, _ := yaml.Marshal(val)
			return string(data)
		},
		"tpl": func(name string, data interface{}) string {
			tmpl, _ := template.New(name).Parse(name)
			var result strings.Builder
			if err := tmpl.Execute(&result, data); err != nil {
				logger.LogError("failed to execute tpl function inside template: %v", err)
				return ""
			}
			return result.String()
		},
	}

	t, err := template.New("tmpl").Funcs(funcMap).Funcs(sprig.GenericFuncMap()).Parse(tmpl)
	if err != nil {
		return "", logger.LogError("Failed to parse template: %v", err)
	}

	var b strings.Builder
	if err := t.Execute(&b, params); err != nil {
		return "", logger.LogError("Failed to execute template with given parameters: %v", err)
	}
	return b.String(), nil
}
