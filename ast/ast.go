package ast

import (
	"fmt"
	"kat/token"
	"strings"
)

const Pad = "   "

type Node interface {
	Node(depth int) string
	String() string
	HasChild() bool
}

// #######################################################
// ##################### Node Program ####################
// #######################################################
type NodeProgram struct {
	Body []Node
}

func (np NodeProgram) Node(depth int) string {
	sb := strings.Builder{}
	sb.WriteString("NodeProgram")

	padding := strings.Repeat(Pad, depth)
	separator := fmt.Sprintf("%s├── ", padding)

	for i, node := range np.Body {
		if i == len(np.Body)-1 {
			separator = fmt.Sprintf("%s└── ", padding)
		}

		sb.WriteString(fmt.Sprintf("\n%s%s", separator, node.Node(depth+1)))
	}

	return sb.String()
}

func (np NodeProgram) String() string {
	return "NodeProgram"
}

func (np NodeProgram) HasChild() bool {
	return true
}

// #######################################################
// ##################### Node Integer ####################
// #######################################################
type NodeInteger struct {
	Token token.Token
	Value int64
}

func (ni NodeInteger) Node(depth int) string {
	return fmt.Sprintf("NodeInteger (%d, %d)", ni.Value, depth)
}

func (ni NodeInteger) String() string {
	return "NodeInteger"
}

func (ni NodeInteger) HasChild() bool {
	return false
}

// #######################################################
// ##################### Node Operator ###################
// #######################################################
type NodeOperator struct {
	Token    token.Token
	Left     Node
	Right    Node
	Operator string
}

func (no NodeOperator) Node(depth int) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("NodeOperator(%s, %d)", no.Operator, depth))

	separator := fmt.Sprintf("│%s", Pad)
	format := strings.Repeat(separator, depth)

	sb.WriteString(fmt.Sprintf("\n%s├── %s", format, no.Left.Node(depth+1)))
	sb.WriteString(fmt.Sprintf("\n%s└── %s", format, no.Right.Node(depth+1)))

	return sb.String()
}

func (no NodeOperator) String() string {
	return "NodeOperator"
}

func (no NodeOperator) HasChild() bool {
	return true
}
