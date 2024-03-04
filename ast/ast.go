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
// ################### Node Boolean #####################
// #######################################################
type NodeBoolean struct {
	Token token.Token
	Value bool
}

func (nb NodeBoolean) Node(indent string) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("NodeBoolean (%t)", nb.Value))
	return sb.String()
}

func (nb NodeBoolean) String() string {
	return "NodeBoolean"
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
	sb.WriteString(fmt.Sprintf("NodeBinary(%s)", nbe.Operator))

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

// #######################################################
// ################## Node Unary Expr ###################
// #######################################################
type NodeUnary struct {
	Token    token.Token
	Operator string
	Right    Node
}

func (np NodeUnary) Node(indent string) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("NodeUnary(%s)", np.Operator))

	_indent := indent
	indent += EmptyPad
	sb.WriteString(fmt.Sprintf("\n%s└── %s", _indent, np.Right.Node(indent)))

	return sb.String()
}

func (np NodeUnary) String() string {
	return "NodeUnary"
}

// #######################################################
// #################### Node Identifier###################
// #######################################################
type NodeIdentifier struct {
	Token token.Token
	Name  string
}

func (ni NodeIdentifier) Node(indent string) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("NodeIdentifier (%s)", ni.Name))

	return sb.String()
}

func (ni NodeIdentifier) String() string {
	return "NodeIdentifier"
}

// #######################################################
// ################### Node Conditional ##################
// #######################################################
type NodeConditionalExpr struct {
	Token     token.Token
	Condition Node
	ThenArm   Node
	ElseArm   Node
}

func (nce NodeConditionalExpr) Node(indent string) string {
	sb := strings.Builder{}

	sb.WriteString("NodeConditionalExpr")

	_indent := indent
	leftIndent := indent + PipePad
	rightIndent := indent + EmptyPad

	sb.WriteString(fmt.Sprintf("\n%s├── %s", _indent, nce.Condition.Node(leftIndent)))
	sb.WriteString(fmt.Sprintf("\n%s└── %s", _indent, "NodeConsequences"))

	_leftIndent := leftIndent + indent + PipePad
	_rightIndent := rightIndent + indent + EmptyPad

	sb.WriteString(fmt.Sprintf("\n%s├── %s", rightIndent, nce.ThenArm.Node(_leftIndent)))
	sb.WriteString(fmt.Sprintf("\n%s└── %s", rightIndent, nce.ElseArm.Node(_rightIndent)))

	return sb.String()
}

func (nce NodeConditionalExpr) String() string {
	return "NodeConditionalExpr"
}
