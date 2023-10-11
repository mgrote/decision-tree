package tree

const (
	CommandType     = "command"
	DecisionType    = "decision"
	DestinationType = "destination"
)

// Command is a step in a branch of a decision tree.
// It is expected that execute function acquire and manipulate data to prepare a next Command, a Decision or a Destination.
type Command struct {
	Name               string      `json:"name"`
	ExpectedInputType  interface{} `json:"inputtype"`
	ExpectedReturnType interface{} `json:"returntype"`
	execute            func(input interface{}) (interface{}, error)
}

// Decision divides a branch of a decision tree in to multiple subbranches.
// It is expected that decide function returns a list of categories names,
// which are used to identify the next nodes to call.
type Decision struct {
	Name              string `json:"name"`
	ExpectedValueType interface{}
	decide            func(input interface{}) ([]string, error)
}

// Destination terminates a branch of a decision tree.
// It is expected that terminate function does an action to fulfill the purpose of the branch.
type Destination struct {
	Name              string `json:"name"`
	ExpectedInputType interface{}
	terminate         func(input interface{}) error
}
