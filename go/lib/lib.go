package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"sigs.k8s.io/yaml"
)

type ReleaseMap map[string]interface{}

type ChartMap map[string]interface{}

type ValuesMap map[string]interface{}

type TemplateData struct {
	Release ReleaseMap
	Chart   ChartMap
	Values  ValuesMap
}

type GetYamlReturnValue struct {
	Yaml    string `json:"yaml"`
	Err     string `json:"err"`
	Warning string `json:"warning"`
}

type GetYamlConfig struct {
	MaxTplRuns int
}

func toJson(returnValue GetYamlReturnValue) string {
	b, err := json.Marshal(returnValue)
	if err != nil {
		return `{"err":"conversion to JSON failed"}`
	}
	return string(b)
}

func GetYaml(templateYaml string, valuesYaml string, config GetYamlConfig) string {
	valuesData := ValuesMap{}
	if err := yaml.Unmarshal([]byte(valuesYaml), &valuesData); err != nil {
		return toJson(GetYamlReturnValue{
			Err: err.Error(),
		})
	}

	releaseData := ReleaseMap{
		"Name":      "example",
		"Namespace": "example",
		"IsUpgrade": false,
		"IsInstall": false,
		"Revision":  1,
		"Service":   "Helm",
	}

	chartData := ChartMap{
		"ApiVersion":   "v2",
		"Name":         "example",
		"Version":      "0.1.0",
		"KubeVersion":  "",
		"Description":  "example",
		"Type":         "application",
		"Keywords":     make([]string, 0),
		"Home":         "",
		"Sources":      make([]string, 0),
		"Dependencies": make([]interface{}, 0),
		"Maintainers":  make([]interface{}, 0),
		"Icon":         "",
		"AppVersion":   "",
		"Deprecated":   false,
		"Annotations":  make(map[string]interface{}),
	}

	templateData := TemplateData{
		Release: releaseData,
		Chart:   chartData,
		Values:  valuesData,
	}

	t := template.New("template")

	var output bytes.Buffer

	funcMap := funcMap()

	/*
	 * Based on:
	 * https://github.com/helm/helm/blob/3d1bc72827e4edef273fb3d8d8ded2a25fa6f39d/pkg/engine/engine.go#L112-L125
	 */
	recursionMaxNums := 5
	includedNames := make(map[string]int)
	funcMap["include"] = func(name string, data interface{}) (interface{}, error) {
		var buf strings.Builder
		if v, ok := includedNames[name]; ok {
			if v > recursionMaxNums {
				return "", errors.New("rendering template has a nested reference name: " + name)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}
		err := t.ExecuteTemplate(&buf, name, data)
		includedNames[name]--
		return buf.String(), err
	}

	// If the template contains `required`, we don't want to return an error which
	// would prevent previewing the entire template. Simply pass the value through.
	funcMap["required"] = func(warn string, val interface{}) (interface{}, error) {
		return val, nil
	}

	// If the template contains `fail`, we don't want to return an error which
	// would prevent previewing the entire template. Return an empty string.
	funcMap["fail"] = func(val interface{}) (interface{}, error) {
		return "", nil
	}

	/*
	 * Based on:
	 * https://github.com/helm/helm/blob/3d1bc72827e4edef273fb3d8d8ded2a25fa6f39d/pkg/engine/engine.go#L128-L152
	 */
	// Add tpl function for rendering templates within templates
	tplRuns := 0
	maxTplRuns := config.MaxTplRuns
	funcMap["tpl"] = func(tpl string, vals interface{}) (interface{}, error) {
		// Simple check to prevent infinite loops
		tplRuns += 1
		if tplRuns >= maxTplRuns {
			return "", errors.New("tpl has been called " + strconv.Itoa(tplRuns) + " times, aborting to prevent infinite loops")
		}

		// Create a new template
		newTemplate := template.New("tpl")

		// Add the function map to the template
		newTemplate.Funcs(funcMap)

		// Parse the template string
		parsed, err := newTemplate.Parse(tpl)
		if err != nil {
			return "", err
		}

		// Execute the template with the provided values
		var buf bytes.Buffer
		if err := parsed.Execute(&buf, vals); err != nil {
			return "", err
		}

		return buf.String(), nil
	}

	t.Funcs(funcMap)

	t, err := t.Parse(templateYaml)
	if err != nil {
		return toJson(GetYamlReturnValue{
			Err: err.Error(),
		})
	}

	if err := t.Execute(&output, templateData); err != nil {
		return toJson(GetYamlReturnValue{
			Err: err.Error(),
		})
	}

	outputYaml := output.String()

	lintValues := ValuesMap{}
	if err := yaml.Unmarshal([]byte(outputYaml), &lintValues); err != nil {
		fmt.Printf("%e", err)
		return toJson(GetYamlReturnValue{
			Warning: err.Error(),
			Yaml:    outputYaml,
		})
	}

	return toJson(GetYamlReturnValue{
		Yaml: outputYaml,
	})
}
