package ast

import (
	"fmt"
	"kat/token"
	"strings"
)

const EmptyPad = "    "
const PipePad = "│   "

type Node interface {
	Node(indent string) string
	String() string
}

// #######################################################
// ##################### Node Program ####################
// #######################################################
type NodeProgram struct {
	Body []Node
}

func (np NodeProgram) Node(indent string) string {
	sb := strings.Builder{}
	sb.WriteString("NodeProgram")

	separator := fmt.Sprintf("%s├── ", indent)

	for i, stmt := range np.Body {
		_indent := indent

		if i == len(np.Body)-1 {
			separator = fmt.Sprintf("%s└── ", _indent)
			_indent += EmptyPad
		} else {
			_indent += PipePad
		}

		sb.WriteString(fmt.Sprintf("\n%s%s", separator, stmt.Node(_indent)))
	}

	return sb.String()
}

func (np NodeProgram) String() string {
	return "NodeProgram"
}

// #######################################################
// ##################### Node Integer ####################
// #######################################################
type NodeInteger struct {
	Token token.Token
	Value int64
}

func (ni NodeInteger) Node(indent string) string {
	return fmt.Sprintf("NodeInteger (%d)", ni.Value)
}

func (ni NodeInteger) String() string {
	return "NodeInteger"
}

// #######################################################
// ##################### Node Operator ###################
// #######################################################
type NodeBinaryExpr struct {
	Token    token.Token
	Left     Node
	Right    Node
	Operator string
}

func (nbe NodeBinaryExpr) Node(indent string) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("NodeBinaryExpr(%s)", nbe.Operator))

	_indent := indent
	leftIndent := indent + PipePad
	rightIndent := indent + EmptyPad

	sb.WriteString(fmt.Sprintf("\n%s├── %s", _indent, nbe.Left.Node(leftIndent)))
	sb.WriteString(fmt.Sprintf("\n%s└── %s", _indent, nbe.Right.Node(rightIndent)))

	return sb.String()
}

func (nbe NodeBinaryExpr) String() string {
	return "NodeBinaryExpr"
}
