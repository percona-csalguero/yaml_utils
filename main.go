package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v2"
)

type cli struct {
	BaseYaml     string   `kong:"name='base-yaml',required"`
	Merge        []string `kong:"name='merge'"`
	SetPath      []string `kong:"name='set-path'"`
	OutputFormat string   `kong:"name=output-format,short='o',enum='json,yaml',default='yaml'"`
}

func main() {
	var opts cli
	parser := kong.Must(&opts,
		kong.Name("poc"),
		kong.Description("POC to merge yaml definitions"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	parser.Parse(os.Args[1:])

	var data map[string]interface{}

	opts.Merge = append([]string{opts.BaseYaml}, opts.Merge...)

	for _, toMerge := range opts.Merge {
		file, err := os.ReadFile(toMerge)
		if err != nil {
			continue
		}

		var mergedata map[string]interface{}

		err = unmarshalTemplate(file, &mergedata)
		if err != nil {
			panic(err)
		}

		data = mergeMaps(data, mergedata)
	}

	for _, pathVal := range opts.SetPath {
		parts := strings.Split(pathVal, "=")
		if len(parts) != 2 {
			panic(errors.New("invalid path/set value"))
		}

		path := strings.Split(parts[0], ".")
		val := parts[1]
		walkToAndReplace(data, path, val)
	}

	var out []byte
	var err error

	switch opts.OutputFormat {
	case "json":
		out, err = json.Marshal(data)
	default:
		out, err = yaml.Marshal(data)
	}
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func walkToAndReplace(m map[string]interface{}, path []string, value interface{}) {
	val, ok := m[path[0]]
	if !ok {
		return
	}

	if len(path) > 1 {
		switch val.(type) {
		case map[string]interface{}:
			walkToAndReplace(m[path[0]].(map[string]interface{}), path[1:], value)
			return
		default:
			return
		}
	}

	m[path[0]] = value
}

func unmarshalTemplate(body []byte, out interface{}) error {
	var yamlObj interface{}

	err := yaml.Unmarshal(body, &yamlObj)
	if err != nil {
		return err
	}

	yamlObj = convert(yamlObj)

	jsonData, err := json.Marshal(yamlObj)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, out)
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := make(map[string]interface{})
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
