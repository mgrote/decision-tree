package tree

import (
	"fmt"
	"github.com/mgrote/meshed/commonmodels"
	"github.com/mgrote/meshed/mesh"
	"log"
	"reflect"
)

func DestinationNodeType() mesh.NodeType {
	return mesh.NewNodeType([]string{commonmodels.CategoryType}, "destination")
}

func init() {
	log.Println("user init called")
	mesh.RegisterTypeConverter("user",
		func() *mesh.Node {
			node := mesh.NewNodeWithContent(DestinationNodeType(), Destination{})
			return &node
		})
	mesh.RegisterContentConverter(DestinationType, GetDestinationFromMap)
}

// NewDestinationNode creates a new destination node
func NewDestinationNode(title string, execFunction func(interface{}) error, inputType, returnType interface{}) (mesh.Node, error) {
	destination := Destination{
		Name:              title,
		terminate:         execFunction,
		ExpectedInputType: inputType,
	}
	node := mesh.NewNodeWithContent(DestinationNodeType(), destination)
	err := node.Save()
	if err != nil {
		return nil, fmt.Errorf("could not save node: %v", err)
	}
	return node, nil
}

// Terminate executes the terminal command with the input.
func Terminate(m mesh.Node, input interface{}) error {
	destination, ok := m.GetContent().(Destination)
	if !ok {
		return fmt.Errorf("could not convert content from %v to destination", m)
	}
	if reflect.TypeOf(input) != destination.ExpectedInputType {
		return fmt.Errorf("input type %v does not match expected input type %v", reflect.TypeOf(input), destination.ExpectedInputType)
	}
	if err := destination.terminate(input); err != nil {
		return fmt.Errorf("could not execute destination: %v", err)
	}
	return nil
}

func GetDestinationFromMap(content map[string]interface{}) interface{} {
	command := Destination{}
	if name, ok := content["name"].(string); ok {
		command.Name = name
	}
	if terminate, ok := content["terminate"].(func(interface{}) error); ok {
		command.terminate = terminate
	}
	if inputType, ok := content["inputtype"].(interface{}); ok {
		command.ExpectedInputType = inputType
	}
	return command
}
