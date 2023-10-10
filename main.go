package main

import (
	"flag"
	"fmt"
	"github.com/mgrote/decision-tree/tree/treemodels/command"
	"github.com/mgrote/meshed/mesh"
	"os"
	"os/exec"
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

	//cmd := exec.Command("ls", "-l")
	////err := cmd.Run()
	////if err != nil {
	////	fmt.Println(err)
	////	os.Exit(1)
	////}
	//out, err := cmd.Output()
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Print(string(out))

	evaluateForMain, err := prepareCommands()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = command.ExecuteCommand(evaluateForMain, ".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func prepareCommands() (mesh.Node, error) {
	commandFunction1 := func(input interface{}) (interface{}, error) {
		cmd := exec.Command("ls", "-l")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		return string(out), nil
	}

	commandFunction2 := func(input interface{}) (interface{}, error) {
		evalString := input.(string)
		out := "false"
		if strings.Contains(evalString, "main.go") {
			out = "true"
		}
		fmt.Println("found main.go", out)
		return out, nil
	}

	commandNode1, err := command.NewCommandNode("ls", commandFunction1, reflect.TypeOf(""), reflect.TypeOf(""))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	commandNode2, err := command.NewCommandNode("has_main", commandFunction2, reflect.TypeOf(""), reflect.TypeOf(""))
	if err != nil {
		return nil, err
	}

	if err = commandNode1.AddChild(commandNode2); err != nil {
		return nil, err
	}

	return commandNode1, nil
}
