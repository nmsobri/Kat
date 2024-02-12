package ast

import (
	"fmt"
	"kat/token"
	"strings"
)

type Node interface {
	String() string
}

type NodeProgram struct {
	Body []Node
}

func (np NodeProgram) String() string {
	sb := strings.Builder{}

	sb.WriteString("NodeProgram{[\n")

	for _, node := range np.Body {
		sb.WriteString("   " + node.String())
		sb.WriteString(",\n")
	}

	sb.WriteString("]}")

	return sb.String()
}

// #######################################################

type NodeInteger struct {
	Token token.Token
	Value int64
}

func (ni NodeInteger) String() string {
	sb := strings.Builder{}

	sb.WriteString("NodeInteger{\n")
	sb.WriteString("   " + fmt.Sprintf("Value: %d", ni.Value))
	sb.WriteString("\n}")

	return sb.String()
}

// #######################################################
