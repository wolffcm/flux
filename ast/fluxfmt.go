package ast

import (
	"fmt"
	"strings"
	"time"
)

const maxLineLength = 80

func (n *Program) String() string {
	ss := []string{}
	for _, s := range n.Body {
		ss = append(ss, s.String())
	}

	return strings.Join(ss, "\n\n")
}

func (n *BlockStatement) String() string {
	ss := []string{}
	for _, s := range n.Body {
		ss = append(ss, s.String())
	}

	return strings.Join(ss, "\n\n")
}

func (n *ExpressionStatement) String() string {
	return n.Expression.String()
}

func (n *ReturnStatement) String() string {
	return n.Argument.String()
}

func (n *OptionStatement) String() string {
	return n.Declaration.String()
}

func (n *VariableDeclaration) String() string {
	ss := []string{}
	for _, s := range n.Declarations {
		ss = append(ss, s.String())
	}

	return strings.Join(ss, "\n")
}

func (n *VariableDeclarator) String() string {
	return n.ID.String() + " = " + n.Init.String()
}

func (n *ArrayExpression) String() string {
	ss := []string{}
	for _, s := range n.Elements {
		ss = append(ss, s.String())
	}

	i := len(ss) - 1

	return "[" + strings.Join(ss[0:i], ", ") + ss[i] + "]"
}

func (n *ArrowFunctionExpression) String() string {
	ss := []string{}
	for _, s := range n.Params {
		ss = append(ss, s.String())
	}

	i := len(ss) - 1
	lhs := "(" + strings.Join(ss[0:i], ", ") + ss[i] + ")"
	rhs := n.Body.String()

	return lhs + " => " + rhs
}

func (n *BinaryExpression) String() string {
	return n.Left.String() + " " + n.Operator.String() + " " + n.Right.String()
}

func (n *CallExpression) String() string {
	ss := []string{}
	for _, s := range n.Arguments {
		ss = append(ss, s.String())
	}
	i := len(ss) - 1
	args := strings.Join(ss[0:i], ", ") + ss[i]

	return n.Callee.String() + "(" + args + ")"
}

func (n *ConditionalExpression) String() string {
	// TODO
	return "if " + n.Test.String() + " then " + n.Consequent.String() + " else " + n.Alternate.String()
}

func (n *LogicalExpression) String() string {
	return n.Left.String() + n.Operator.String() + n.Right.String()
}

func (n *MemberExpression) String() string {
	// TODO
	return n.Object.String() + n.Property.String()
}

func (n *PipeExpression) String() string {
	return n.Argument.String() + "\n  |> " + n.Call.String()
}

func (n *ObjectExpression) String() string {
	ss := []string{}
	for _, s := range n.Properties {
		ss = append(ss, s.String())
	}

	i := len(ss) - 1

	return "{" + strings.Join(ss[0:i], ", ") + ss[i] + "}"
}

func (n *UnaryExpression) String() string {
	return n.Operator.String() + n.Argument.String()
}

func (n *Property) String() string {
	if n.Value != nil {
		return n.Key.String() + ": " + n.Value.String()
	}

	return n.Key.String()
}

func (n *Identifier) String() string {
	return n.Name
}

func (n *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", n.Value)
}

func (n *DateTimeLiteral) String() string {
	return n.Value.Format(time.RFC3339)
}

func (n *Duration) String() string {
	return fmt.Sprintf("%d%s", n.Magnitude, n.Unit)
}

func (n *DurationLiteral) String() string {
	ss := []string{}
	for _, s := range n.Values {
		ss = append(ss, s.String())
	}

	return strings.Join(ss, "")
}

func (n *FloatLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (n *PipeLiteral) String() string {
	return "|>"
}

func (n *RegexpLiteral) String() string {
	return `/` + n.Value.String() + `/`
}

func (n *StringLiteral) String() string {
	return `"` + n.Value + `"`
}

func (n *UnsignedIntegerLiteral) String() string {
	return fmt.Sprintf("%d", n.Value)
}
