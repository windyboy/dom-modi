package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/integrii/flaggy"
	"gopkg.in/yaml.v2"
)

// Declare variables and their defaults
var config = "cfg.yml"
var sourceFileName = "domain.xml"
var sourceFile string
var targetFile string
var verbose bool
var targetDir = "modify/"
var overWrite bool
var rules []Expression

//Expression config
type Expression struct {
	Name       string
	Enable     bool
	Expression string
	Replace    string
}

// LoadConfig read config file content
func LoadConfig(fileLoc string) (string, error) {
	content, err := ioutil.ReadFile(config)
	return string(content), err
}

// MarshalConfig config yaml file
func MarshalConfig(config string) ([]Expression, error) {
	r := []Expression{}
	err := yaml.Unmarshal([]byte(config), &r)
	return r, err
}

// FindExpressionByName get reglular expression from config by name
func FindExpressionByName(rule []Expression, name string) Expression {
	for _, exp := range rule {
		if name == exp.Name {
			return exp
		}
	}
	return Expression{}
}

func matchAndReplace(source string, rule Expression) string {
	if rule.Enable {
		expression, err := regexp.Compile(rule.Expression)
		if err != nil {
			fmt.Fprintf(os.Stderr, "expression compile error %v - expression %s\n", err, expression)
			return source
		} else if expression.MatchString(sourceFile) {
			fmt.Printf("find match with rule [%s]", rule.Name)
			return expression.ReplaceAllString(source, rule.Replace)
		}
	}
	return source
}

func main() {

	// set a description, name, and version for our parser
	p := flaggy.NewParser("dom-modi")
	p.Description = "modify the kvm domain xml exported from centos kvm."
	p.Version = "0.0.1"
	// display some before and after text for all help outputs
	// p.AdditionalHelpPrepend = "I hope you like this program!"
	// p.AdditionalHelpAppend = "This command has no warranty."

	// add a positional value at position 1
	p.AddPositionalValue(&sourceFileName, "source", 1, true, "kvm domain xml file location")

	// create a subcommand at position 2
	// you don't have to finish the subcommand before adding it to the parser
	modifyCmd := flaggy.NewSubcommand("modify")
	modifyCmd.Description = "modify the source file"

	modifyCmd.String(&config, "c", "config", "rules of modification")
	modifyCmd.String(&targetDir, "t", "target", "modified file save location")

	p.AttachSubcommand(modifyCmd, 2)

	// add a bool flag to the root command
	p.Bool(&verbose, "v", "verbose", "detail information")

	p.Bool(&overWrite, "o", "overwrite", "overwrite the target file")

	p.Parse()

	//load source file
	sourceFileBytes, err := ioutil.ReadFile(sourceFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "source file %s error %v\n", sourceFileName, err)
		os.Exit(1)
	}
	sourceFile = string(sourceFileBytes)
	if verbose {
		fmt.Printf("source %s loaded %d", sourceFileName, len(sourceFile))
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "target directory %s is not exists ", targetDir)
		os.Exit(1)
	}

	// content, err := ioutil.ReadFile(config)
	content, err := LoadConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// r, err := mashalConfig(string(content))
	rules, err := MarshalConfig(content)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// fmt.Println(len(rules))
	targetFile = sourceFile
	for _, rule := range rules {
		fmt.Printf("[%s] > enable[%t] \n", rule.Name, rule.Enable)
		targetFile = matchAndReplace(targetFile, rule)
	}
	fullTargetFile := targetDir + "/" + sourceFileName
	if err := ioutil.WriteFile(fullTargetFile, []byte(targetFile), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "wirte modified content failed %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf(" file [%s] wirte %d bytes", fullTargetFile, len(targetFile))
	}
	// fmt.Println(sourceFile, config, targetDir, verbose)

}
