/*
 *	Funcscript Syntax
 */

package funcscript

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/lexer"
)

type Script struct {
	Commands []*Command `parser:"(@@ ';')*"`
}

type Command struct {
	Pos  lexer.Position
	Name string       `parser:"@Ident '('"`
	Args []Expression `parser:"(@@ (',' @@)*)? ')'"`
}

type Expression struct {
	Pos          lexer.Position
	IntValue     *int     `parser:"  @Int"`
	DoubleValue  *float64 `parser:"| @Float"`
	BoolValueStr *string  `parser:"| @('true' | 'false')"`
	StringValue  *string  `parser:"| @String"`
	CommandValue *Command `parser:"| @@"`
}

func (c *Command) String() string {
	sb := strings.Builder{}
	sb.WriteString(c.Name)
	if len(c.Args) == 0 {
		sb.WriteString("()")
		return sb.String()
	}

	sb.WriteRune('(')
	sb.WriteString(c.Args[0].String())
	for i := 1; i < len(c.Args); i++ {
		sb.WriteString(", ")
		sb.WriteString(c.Args[i].String())
	}
	sb.WriteRune(')')
	return sb.String()
}

func (e *Expression) BoolValue() bool {
	return e.BoolValueStr != nil && *e.BoolValueStr == "true"
}

func (e *Expression) String() string {
	if e.IntValue != nil {
		return fmt.Sprintf("%d", *e.IntValue)
	}
	if e.DoubleValue != nil {
		return fmt.Sprintf("%g", *e.DoubleValue)
	}
	if e.BoolValueStr != nil {
		return *e.BoolValueStr
	}
	if e.StringValue != nil {
		return *e.StringValue
	}
	if e.CommandValue != nil {
		return e.CommandValue.String()
	}
	return ""
}
