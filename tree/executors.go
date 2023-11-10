package tree

import (
	"fmt"
	"github.com/mgrote/meshed/mesh"
)

func ExecuteNodes(nodes []mesh.Node, input interface{}, executorName string) error {
	// aggregate errors
	var aggregated error
	var err error
	for _, node := range nodes {
		switch node.GetTypeName() {
		case CommandType:
			if err = ExecuteCommand(node, input); err != nil {
				aggregated = fmt.Errorf("node %s -> %w", executorName, err)
			}
		case DecisionType:
			if err = Decide(node, input); err != nil {
				aggregated = fmt.Errorf("node %s -> %w", executorName, err)
			}
		case DestinationType:
			if err = Terminate(node, input); err != nil {
				aggregated = fmt.Errorf("node %s -> %w", executorName, err)
			}
		}
	}
	return aggregated
}
