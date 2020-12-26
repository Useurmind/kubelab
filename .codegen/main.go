package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

type CodegenConfig struct {
	ModelFolder string `yaml:"modelFolder"`
	TemplateFolder string `yaml:"templateFolder"`
	TemplateName string `yaml:"templateName"`
	PocoFolders map[string]string `yaml:"pocoFolders"`
}

type PocoSet struct {
	GoNamespace string `yaml:"goNamespace"`

	PocoTypes []PocoDefinition `yaml:"pocoTypes"`
}

type PocoDefinition struct {
	PocoName string `yaml:"pocoName"`
	Description string `yaml:"description"`
	Properties []PropertyDefinition `yaml:"properties"`
}

type PropertyDefinition struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Description string `yaml:"description"`
	TypeMap map[string]string `yaml:"typeMap"`
}

func (p PropertyDefinition) GetType(key string) string {
	if p.Type != "" {
		return p.Type
	}

	return p.TypeMap[key]
}

// this code genera
func main() {	
	var configFile string
	if len(os.Args) < 2 {
		configFile = ".codegen-config.yml"
	} else {
		configFile = os.Args[1]
	}

	config, err := ReadCodegenConfig(configFile)
	if err != nil {
		handleError(fmt.Sprintf("Could not read codegen config file %s", configFile), err)
	}

	err = GenerateCode(config)
	if err != nil {
		handleError("Could not generate code", err)
	}
}

func handleError(msg string, err error) {
	fmt.Printf("ERROR: %s - %v\n", msg, err)
	os.Exit(1)
}

func GenerateCode(config CodegenConfig) error {
	modelFileInfos, err := ioutil.ReadDir(config.ModelFolder)
	if err != nil {
		return fmt.Errorf("Could not read model directory %s - %v", config.ModelFolder, err)
	}

	for lang, pocoFolder := range config.PocoFolders {
		templateFile := config.TemplateFolder + "/" + config.TemplateName + "." + lang

		for _, modelFile := range modelFileInfos {
			specFile := config.ModelFolder + "/" + modelFile.Name()
			outputFile := pocoFolder + "/" + strings.ReplaceAll(modelFile.Name(), filepath.Ext(modelFile.Name()), "." + lang)

			spec, err := ReadSpecFile(specFile)
			if err != nil {
				return err
			}

			outputDir := path.Dir(outputFile)
			err = os.MkdirAll(outputDir, 777)
			if err != nil {
				return err
			}

			err = FillTemplate(templateFile, spec, outputFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadCodegenConfig(configFile string) (config CodegenConfig, err error) {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return  config, err
	}

	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return config, fmt.Errorf("Failed to parse yaml values: %v", err)
	}

	fmt.Printf("Codegen config is:\n%#v\n", config)

	return config, nil
}

func ReadSpecFile(specFile string) (pocoSet PocoSet, err error) {
	bytes, err := ioutil.ReadFile(specFile)
	if err != nil {
		return  pocoSet, err
	}

	err = yaml.Unmarshal(bytes, &pocoSet)
	if err != nil {
		return pocoSet, fmt.Errorf("Failed to parse yaml values: %v", err)
	}

	fmt.Printf("Parse poco set is:\n%#v\n", pocoSet)

	return pocoSet, nil
}

// FillTemplate fills the file at the templateFilePath with values taken from the valuesFilePath.
// It prints the filled template to the outputFilePath.
// You can specify additional values that are not contained in the valuesFilePath file.
// The valuesFilePath must point to a valid yaml file.
func FillTemplate(templateFilePath string, pocoSet PocoSet, outputFilePath string) error {
	slashTemplatePath := filepath.ToSlash(templateFilePath)

	fmt.Printf("Filling %s to %s\n", templateFilePath, outputFilePath)

	outputFile, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tpl, err := template.New(path.Base(slashTemplatePath)).Funcs(sprig.TxtFuncMap()).ParseFiles(slashTemplatePath)
	if err != nil {
		return fmt.Errorf("Error parsing template(s): %v", err)
	}

	values := map[string]interface{} {
		"Spec": pocoSet,
	}

	err = tpl.Execute(outputFile, values)
	if err != nil {
		return fmt.Errorf("Failed to execute template: %v", err)
	}
	return nil
}