package ast

import (
	"fmt"
	"kat/token"
	"strings"
)

const pad = "    "

type Node interface {
	String(depth int, isLast bool) string
}

type NodeProgram struct {
	Body []Node
}

func (np NodeProgram) String(depth int, isLast bool) string {
	sb := strings.Builder{}
	padding := strings.Repeat(pad, depth)
	indent := fmt.Sprintf("%s├── ", padding)

	sb.WriteString("NodeProgram")

	for i, node := range np.Body {
		if i == len(np.Body)-1 {
			indent = fmt.Sprintf("%s└── ", padding)
			isLast = true
		}

		sb.WriteString("\n")
		sb.WriteString(indent)
		sb.WriteString(node.String(depth+1, isLast))
	}

	return sb.String()
}

// #######################################################

type NodeInteger struct {
	Token token.Token
	Value int64
}

func (ni NodeInteger) String(depth int, isLast bool) string {
	return fmt.Sprintf("NodeInteger (%d)", ni.Value)
}

// #######################################################

type NodeInfix struct {
	Token    token.Token
	Left     Node
	Right    Node
	Operator string
}

func (ni NodeInfix) String(depth int, isLast bool) string {
	sb := strings.Builder{}
	padding := strings.Repeat(pad, depth)
	format := "│"

	if isLast {
		format = ""
	}

	sb.WriteString(fmt.Sprintf("NodeInfix(%s)", ni.Operator))
	sb.WriteString(fmt.Sprintf("\n%s%s├── ", format, padding))
	sb.WriteString(ni.Left.String(depth+1, isLast))
	sb.WriteString(fmt.Sprintf("\n%s%s└── ", format, padding))
	sb.WriteString(ni.Right.String(depth+1, isLast))

	return sb.String()

}
