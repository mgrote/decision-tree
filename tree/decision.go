package tree

import (
	"fmt"
	"github.com/mgrote/meshed/commonmodels"
	"github.com/mgrote/meshed/commonmodels/categories"
	"github.com/mgrote/meshed/mesh"
	"log"
	"reflect"
	"slices"
)

func DecisionNodeType() mesh.NodeType {
	return mesh.NewNodeType([]string{commonmodels.CategoryType}, "decision")
}

func init() {
	log.Println("user init called")
	mesh.RegisterTypeConverter("user",
		func() *mesh.Node {
			node := mesh.NewNodeWithContent(DecisionNodeType(), Decision{})
			return &node
		})
	mesh.RegisterContentConverter(DecisionType, GetDecisionFromMap)
}

// NewDecisionNode creates a new decision node.
func NewDecisionNode(title string, execFunction func(interface{}) ([]string, error), valueType, returnType interface{}) (mesh.Node, error) {
	decision := Decision{
		Name:              title,
		decide:            execFunction,
		ExpectedValueType: valueType,
	}
	node := mesh.NewNodeWithContent(DecisionNodeType(), decision)
	err := node.Save()
	if err != nil {
		return nil, fmt.Errorf("could not save node: %v", err)
	}
	return node, nil
}

// Decide executes the decision command with the input.
// The input is expected to be unchanged.
// The outcome of a decision is a list of categories names.
// These categories are used to identify the next nodes to call.
func Decide(m mesh.Node, input interface{}) error {
	decision, ok := m.GetContent().(Decision)
	if !ok {
		return fmt.Errorf("could not convert content from %v to decision", m)
	}

	if reflect.TypeOf(input) != decision.ExpectedValueType {
		return fmt.Errorf("input type %v does not match expected value type %v", reflect.TypeOf(input), decision.ExpectedValueType)
	}

	decisions, err := decision.decide(input)
	if err != nil {
		return fmt.Errorf("could not make decisions: %v", err)
	}

	// Search nodes with matching categories from decision and execute them.
	matchingNodes := getNodesByCategoryNames(m, decisions)
	return ExecuteNodes(matchingNodes, input, decision.Name)
}

func GetDecisionFromMap(content map[string]interface{}) interface{} {
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

func getNodesByCategoryNames(m mesh.Node, names []string) []mesh.Node {
	var nodes []mesh.Node
	for _, node := range m.GetChildrenIn(CommandType, DecisionType, DestinationType) {
		catNodes := node.GetNodes(commonmodels.CategoryType)
		for _, cat := range catNodes {
			if slices.Contains(names, categories.GetCategory(cat).Name) {
				nodes = append(nodes, node)
			}
		}
	}
	return nodes
}
