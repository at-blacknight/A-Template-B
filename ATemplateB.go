package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

// Version variable will hold the version info at build time
var version string

// TemplateContext will wrap DynamicConfig to expose it as .Values in the template.
type TemplateContext struct {
	Values map[string]interface{} // Root context (DynamicConfig)
}

// DynamicConfig will hold values from JSON/YAML
var DynamicConfig map[string]interface{}

func main() {
	// Command-line arguments
	versionFlag := flag.Bool("version", false, "Display the application version")
	templateFile := flag.String("template", "", "Path to the Nginx template file")
	valuesFile := flag.String("values", "", "Path to the JSON or YAML file with values")
	outputFile := flag.String("output", "nginx.conf", "Path to output file")
	flag.Parse()

	//always display version in use
	fmt.Println("Application Version:", version)
	// If --version flag is passed exit immediately
	if *versionFlag {
		os.Exit(0)
	}

	// Validate inputs
	if *templateFile == "" {
		fmt.Println("Usage: ATemplateB --template=<template>.tmpl[.<extension>] [--values=<values>.yaml/json] --output=<output>.<extension>")
		os.Exit(1)
	}

	// Load and parse the values file
	if *valuesFile != "" {
		data, err := os.ReadFile(*valuesFile)
		if err != nil {
			panic(err)
		}

		// Determine if the values file is JSON or YAML
		if json.Valid(data) {
			err = json.Unmarshal(data, &DynamicConfig)
		} else {
			err = yaml.Unmarshal(data, &DynamicConfig)
		}
		if err != nil {
			panic(err)
		}
	} else {
		// No values file provided, initialize empty map to avoid nil pointer in template
		fmt.Println("You have not provided a values file, an empty map will be used. if this was not intentional please review the usage:")
		fmt.Println("ATemplateB --template=<template>.tmpl[.<extension>] [--values=<values>.yaml/json] --output=<output>.<extension>")
		DynamicConfig = make(map[string]interface{})
	}

	// Wrap DynamicConfig in TemplateContext
	context := TemplateContext{
		Values: DynamicConfig,
	}

	// Add getFile support
	funcMap := sprig.FuncMap()
	funcMap["getFile"] = func(path string) string {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Sprintf("# ERROR: could not read %s: %v", path, err)
		}
		return string(content)
	}
	// getYaml: parse & pretty format YAML
	funcMap["getYaml"] = func(path string) string {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Sprintf("# ERROR: could not read %s: %v", path, err)
		}

		var node yaml.Node
		err = yaml.Unmarshal(data, &node)
		if err != nil {
			return fmt.Sprintf("# ERROR: could not parse YAML %s: %v", path, err)
		}

		formatted, err := yaml.Marshal(&node)
		if err != nil {
			return fmt.Sprintf("# ERROR: could not format YAML %s: %v", path, err)
		}

		return string(formatted)
	}

	// Extract the directory from the provided template file path
	templateDir := filepath.Dir(*templateFile)
	templateName := filepath.Base(*templateFile) // Get just the file name
	// Parse all templates in the same directory
	tmpl, err := template.New(templateName).Funcs(funcMap).ParseGlob(filepath.Join(templateDir, "*.tmpl*"))
	if err != nil {
		panic(err)
	}

	// Create and write the final config file
	output, err := os.Create(*outputFile)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	// Execute the originally provided template
	err = tmpl.ExecuteTemplate(output, templateName, context)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated config saved to:", *outputFile)
}
