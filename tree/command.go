package tree

import (
	"fmt"
	"github.com/mgrote/meshed/commonmodels"
	"github.com/mgrote/meshed/mesh"
	"log"
	"reflect"
)

func CommandNodeType() mesh.NodeType {
	return mesh.NewNodeType([]string{CommandType, DecisionType, DestinationType, commonmodels.CategoryType}, "command")
}

func init() {
	log.Println("user init called")
	mesh.RegisterTypeConverter("user",
		func() *mesh.Node {
			node := mesh.NewNodeWithContent(CommandNodeType(), Command{})
			return &node
		})
	mesh.RegisterContentConverter(CommandType, GetCommandFromMap)
}

// NewCommandNode creates a new command node
func NewCommandNode(title string, execFunction func(interface{}) (interface{}, error), inputType, returnType interface{}) (mesh.Node, error) {
	command := Command{
		Name:               title,
		execute:            execFunction,
		ExpectedInputType:  inputType,
		ExpectedReturnType: returnType,
	}
	node := mesh.NewNodeWithContent(CommandNodeType(), command)
	err := node.Save()
	if err != nil {
		return nil, fmt.Errorf("could not save node: %v", err)
	}
	return node, nil
}

// ExecuteCommand executes the command and checks if the input and return types match the expected types.
func ExecuteCommand(m mesh.Node, input interface{}) error {
	command, ok := m.GetContent().(Command)
	if !ok {
		return fmt.Errorf("%s.%s '%s': could not convert content from %v to command", m.GetTypeName(), m.GetID(), command.Name, m)
	}

	if isSliceOrArray(input) && reflect.TypeOf(input) != command.ExpectedInputType {
		return aggregatedExecuteCommand(m, command, input)
	}

	return executeSingleValueCommand(m, command, input)
}

// aggregatedExecuteCommand is called with a slice or array as input and iterates over the input elements.
func aggregatedExecuteCommand(m mesh.Node, command Command, input interface{}) error {

	if reflect.TypeOf(input).Elem() != command.ExpectedInputType {
		return fmt.Errorf("%s.%s '%s': input type %v does not match expected input type %v", m.GetTypeName(), m.GetID(), command.Name, reflect.TypeOf(input), command.ExpectedInputType)
	}

	var aggregated error
	var err error
	elements := reflect.ValueOf(input)
	debugErrorCounter := 0
	for j := 0; j < elements.Len(); j++ {
		if err = executeSingleValueCommand(m, command, elements.Index(j).Interface()); err != nil {
			aggregated = fmt.Errorf("node %s: value \"%s\" - %d: %w", command.Name, elements.Index(j), debugErrorCounter, err)
		}
		debugErrorCounter++
	}

	return aggregated
}

func executeSingleValueCommand(m mesh.Node, command Command, input interface{}) error {

	if reflect.TypeOf(input) != command.ExpectedInputType {
		return fmt.Errorf("%s.%s '%s': input type %v does not match expected input type %v", m.GetTypeName(), m.GetID(), command.Name, reflect.TypeOf(input), command.ExpectedInputType)
	}

	out, err := command.execute(input)
	if err != nil {
		return fmt.Errorf("%s.%s '%s': could not execute command: %v", m.GetTypeName(), m.GetID(), command.Name, err)
	}

	if reflect.TypeOf(out) != command.ExpectedReturnType {
		return fmt.Errorf("%s.%s '%s': return type %v does not match expected output type %v", m.GetTypeName(), m.GetID(), command.Name, reflect.TypeOf(out), command.ExpectedReturnType)
	}

	return ExecuteNodes(m.GetChildrenIn(CommandType, DecisionType, DestinationType), out, command.Name)
}

func GetCommandFromMap(content map[string]interface{}) interface{} {
	command := Command{}
	if name, ok := content["name"].(string); ok {
		command.Name = name
	}
	if exec, ok := content["exec"].(func(interface{}) (interface{}, error)); ok {
		command.execute = exec
	}
	if inputType, ok := content["inputtype"].(interface{}); ok {
		command.ExpectedInputType = inputType
	}
	if returnType, ok := content["returntype"].(interface{}); ok {
		command.ExpectedReturnType = returnType
	}
	return command
}

func isSliceOrArray(input interface{}) bool {
	return reflect.TypeOf(input).Kind() == reflect.Slice || reflect.TypeOf(input).Kind() == reflect.Array
}
