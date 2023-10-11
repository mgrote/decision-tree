package tree

import (
	"fmt"
	"github.com/mgrote/meshed/mesh"
	"log"
	"reflect"
)

func CommandNodeType() mesh.NodeType {
	return mesh.NewNodeType([]string{CommandType, DecisionType}, "command")
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
		return fmt.Errorf("%s.%s: could not convert content from %v to command", m.GetTypeName(), m.GetID().String(), m)
	}

	if reflect.TypeOf(input) != command.ExpectedInputType {
		return fmt.Errorf("%s.%s: input type %v does not match expected input type %v", m.GetTypeName(), m.GetID().String(), reflect.TypeOf(input), command.ExpectedInputType)
	}

	out, err := command.execute(input)
	if err != nil {
		return fmt.Errorf("%s.%s: could not execute command: %v", m.GetTypeName(), m.GetID().String(), err)
	}

	if reflect.TypeOf(out) != command.ExpectedReturnType {
		return fmt.Errorf("%s.%s: return type %v does not match expected output type %v", m.GetTypeName(), m.GetID().String(), reflect.TypeOf(out), command.ExpectedReturnType)
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
