package treemodels

import (
	"fmt"
	"github.com/mgrote/decision-tree/tree"
	"github.com/mgrote/decision-tree/tree/treemodels/command"
	"github.com/mgrote/decision-tree/tree/treemodels/decision"
	"github.com/mgrote/decision-tree/tree/treemodels/destination"
	"github.com/mgrote/meshed/mesh"
)

//type TreeNode interface {
//	Execute() error
//}

func ExecuteNodes(nodes []mesh.Node, input interface{}, executorName string) error {
	// aggregate errors
	var aggregated error
	var err error
	for _, node := range nodes {
		switch node.GetTypeName() {
		case tree.CommandType:
			if err = command.ExecuteCommand(node, input); err != nil {
				aggregated = fmt.Errorf("%s: %w", executorName, err)
			}
		case tree.DecisionType:
			if err = decision.Decide(node, input); err != nil {
				aggregated = fmt.Errorf("%s: %w", executorName, err)
			}
		case tree.DestinationType:
			if err = destination.Terminate(node, input); err != nil {
				aggregated = fmt.Errorf("%s: %w", executorName, err)
			}
		}
	}
	return aggregated
}
