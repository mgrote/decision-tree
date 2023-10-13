package main

import (
	"flag"
	"fmt"
	"github.com/mgrote/decision-tree/tree"
	"github.com/mgrote/meshed/commonmodels/categories"
	"github.com/mgrote/meshed/mesh"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
)

func main() {

	var pathFlag string
	flag.StringVar(&pathFlag, "inifiles", "./config/mesh.db.properties.ini", "Path to ini files")
	flag.Parse()

	// Init API with default config.
	if err := mesh.InitApiWithConfig(pathFlag); err != nil {
		fmt.Println("init mesh api:", err)
		os.Exit(1)
	}

	// prepare compile main or fail decision tree
	evaluateForMainRoot, err := prepareCommands()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// start the decision tree
	err = tree.ExecuteCommand(evaluateForMainRoot, ".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func prepareCommands() (mesh.Node, error) {

	categoryMainFound := categories.NewNode("contains_main")
	categorySearchMain := categories.NewNode("search_main")
	categoryNotFound := categories.NewNode("not_found")

	commandListCurrentDir := func(input interface{}) (interface{}, error) {
		cmd := exec.Command("ls", "-l")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		return string(out), nil
	}

	decisionMainFound := func(input interface{}) ([]string, error) {
		evalString := input.(string)
		var cats []string
		if strings.Contains(evalString, "main.go") {
			cats = append(cats, categories.GetCategory(categoryMainFound).Name)
		} else if strings.Contains(evalString, "dr") {
			// TODO: have to evaluate each line to start with 'dr'
			cats = append(cats, categories.GetCategory(categorySearchMain).Name)
		} else {
			cats = append(cats, categories.GetCategory(categoryNotFound).Name)
		}
		return cats, nil
	}

	commandFindDirNames := func(input interface{}) (interface{}, error) {
		cmd := exec.Command("find", ".", "-maxdepth", "1", "-type", "d")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		return string(out), nil
	}

	commandChangeDir := func(input interface{}) (interface{}, error) {
		cmd := exec.Command("cd", input.(string))
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		cmd = exec.Command("pwd")
		out, err = cmd.Output()
		if err != nil {
			return nil, err
		}
		return string(out), nil
	}

	completePathCommandToMain := func(input interface{}) (interface{}, error) {
		cmd := exec.Command("pwd")
		out, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("completePathToMain: could not execute pwd: %w", err)
		}
		completePath := filepath.Join(string(out), "main.go")
		fmt.Println("found main.go", out)
		return completePath, nil
	}

	compileMainDestination := func(input interface{}) error {
		cmd := exec.Command("go", "build", input.(string))
		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("compileMainDestination: could not execute go build: %w", err)
		}
		fmt.Println("compiled main.go", out)
		return nil
	}

	nothingFoundDestination := func(input interface{}) error {
		fmt.Println("nothing found")
		return fmt.Errorf("nothing found to compile")
	}

	lsDirNode, err := tree.NewCommandNode("ls", commandListCurrentDir, reflect.TypeOf(""), reflect.TypeOf(""))
	if err != nil {
		handleError(err)
	}

	chDirNode, err := tree.NewCommandNode("cd", commandChangeDir, reflect.TypeOf(""), reflect.TypeOf(""))
	if err != nil {
		handleError(err)
	}

	findDirNode, err := tree.NewCommandNode("find", commandFindDirNames, reflect.TypeOf(""), reflect.TypeOf(""))
	if err != nil {
		handleError(err)
	}
	if err := findDirNode.AddChild(categorySearchMain); err != nil {
		handleError(err)
	}

	completePathNode, err := tree.NewCommandNode("completePathToMain", completePathCommandToMain, reflect.TypeOf(""), reflect.TypeOf(""))
	if err != nil {
		handleError(err)
	}
	if err := completePathNode.AddChild(categoryMainFound); err != nil {
		handleError(err)
	}

	compileMainNode, err := tree.NewDestinationNode("compileMain", compileMainDestination, reflect.TypeOf(""), reflect.TypeOf(""))
	if err != nil {
		handleError(err)
	}

	nothingFoundNode, err := tree.NewDestinationNode("nothingFound", nothingFoundDestination, reflect.TypeOf(""), reflect.TypeOf(""))
	if err != nil {
		handleError(err)
	}
	if err := nothingFoundNode.AddChild(categoryNotFound); err != nil {
		handleError(err)
	}

	decisionNode, err := tree.NewDecisionNode("decision", decisionMainFound, reflect.TypeOf(""), reflect.TypeOf(""))
	if err != nil {
		handleError(err)
	}

	if err = lsDirNode.AddChild(decisionNode); err != nil {
		handleError(err)
	}
	// prepare decision
	if err = decisionNode.AddChild(completePathNode); err != nil {
		handleError(err)
	}
	if err = decisionNode.AddChild(findDirNode); err != nil {
		handleError(err)
	}
	if err = decisionNode.AddChild(nothingFoundNode); err != nil {
		handleError(err)
	}

	// prepare compile main
	if err = completePathNode.AddChild(compileMainNode); err != nil {
		handleError(err)
	}

	// prepare search for main
	if err = findDirNode.AddChild(chDirNode); err != nil {
		handleError(err)
	}
	if err = chDirNode.AddChild(lsDirNode); err != nil {
		handleError(err)
	}

	return lsDirNode, nil
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
