package decision

import (
	"fmt"
	"github.com/mgrote/decision-tree/tree"
	"github.com/mgrote/meshed/commonmodels"
	"github.com/mgrote/meshed/mesh"
	"log"
	"reflect"
)

func NodeType() mesh.NodeType {
	return mesh.NewNodeType([]string{commonmodels.CategoryType}, "decision")
}

type Decision struct {
	// The name of the command.
	Name string `json:"name"`
	// The expected input type.
	ExpectedValueType interface{}
	decide            func(input interface{}) ([]string, error)
}

func init() {
	log.Println("user init called")
	mesh.RegisterTypeConverter("user",
		func() *mesh.Node {
			node := mesh.NewNodeWithContent(NodeType(), Decision{})
			return &node
		})
	mesh.RegisterContentConverter(tree.DecisionType, GetFromMap)
}

// NewDestinationNode creates a new destination node
func NewDestinationNode(title string, execFunction func(interface{}) ([]string, error), valueType, returnType interface{}) (mesh.Node, error) {
	decision := Decision{
		Name:              title,
		decide:            execFunction,
		ExpectedValueType: valueType,
	}
	node := mesh.NewNodeWithContent(NodeType(), decision)
	err := node.Save()
	if err != nil {
		return nil, fmt.Errorf("could not save node: %v", err)
	}
	return node, nil
}

// Decide executes the decision command with the input.
func Decide(m mesh.Node, input interface{}) error {
	destination, ok := m.GetContent().(Decision)
	if !ok {
		return fmt.Errorf("could not convert content from %v to decision", m)
	}
	if reflect.TypeOf(input) != destination.ExpectedValueType {
		return fmt.Errorf("input type %v does not match expected value type %v", reflect.TypeOf(input), destination.ExpectedValueType)
	}
	decisions, err := destination.decide(input)
	if err != nil {
		return fmt.Errorf("could not execute destination: %v", err)
	}
	// search nodes with matching categories from decisions
	log.Println(decisions)
	return nil
}
func GetFromMap(content map[string]interface{}) interface{} {
	command := Decision{}
	if name, ok := content["name"].(string); ok {
		command.Name = name
	}
	if decide, ok := content["decide"].(func(interface{}) ([]string, error)); ok {
		command.decide = decide
	}
	if inputType, ok := content["valuetype"].(interface{}); ok {
		command.ExpectedValueType = inputType
	}
	return command
}
