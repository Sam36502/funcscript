package funcscript

import (
	"fmt"

	"github.com/alecthomas/participle"
)

var g_parser *participle.Parser
var g_funcs map[string]HandlerFunction
var g_vars map[string]Expression

type HandlerFunction func(ctx Context) (*Expression, error)

func Initialise() error {
	var err error

	g_parser, err = participle.Build(&Script{})
	if err != nil {
		return fmt.Errorf("failed to build FuncScript parser: %v", err)
	}

	g_funcs = make(map[string]HandlerFunction)
	for n, f := range g_builtins {
		err = AddFunction(n, f)
		if err != nil {
			return fmt.Errorf("overlap in builtin names: %v", err)
		}
	}

	g_vars = make(map[string]Expression)

	return nil
}

// Registers a function to be usable in a funcscript string
func AddFunction(name string, f HandlerFunction) error {
	if _, exists := g_funcs[name]; exists {
		return fmt.Errorf("function '%s' is already registered", name)
	}
	g_funcs[name] = f
	return nil
}

// Lists all currently registered functions
func GetFunctions() []string {
	fs := []string{}
	for n := range g_funcs {
		fs = append(fs, n)
	}
	return fs
}

// Sets a global variable to be accessed from funcscripts with '_get'
func SetVar(name string, value Expression) {
	g_vars[name] = value
}

// Gets a global variable accessable from funcscripts with '_get'
func GetVar(name string) (Expression, error) {
	v, e := g_vars[name]
	if !e {
		return Expression{}, fmt.Errorf("variable '%s' is not set", name)
	}
	return v, nil
}

// Evaluates a funcscript string
func Eval(script string) (*Expression, error) {
	ast := Script{}
	err := g_parser.ParseString(script, &ast)
	if err != nil {
		return nil, err
	}

	var expr *Expression
	for _, cmd := range ast.Commands {
		expr, err = evalCommand(*cmd)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", cmd.Pos, err)
		}
	}

	return expr, nil
}

func evalCommand(cmd Command) (*Expression, error) {
	ctx := Context{
		FuncName: cmd.Name,
		Args:     cmd.Args,
	}

	// Handle special cases before args are evaluated
	switch cmd.Name {
	case "_if":
		return biIf(ctx)
	case "_while":
		return biWhile(ctx)
	}

	f, exists := g_funcs[cmd.Name]
	if !exists {
		return nil, fmt.Errorf("no function '%s' is registered", cmd.Name)
	}

	// Check recursively for nested commands and eval args
	ctx.Args = make([]Expression, len(cmd.Args))
	for i, arg := range cmd.Args {
		if arg.CommandValue != nil {
			expr, err := evalCommand(*arg.CommandValue)
			if err != nil {
				return nil, err
			}
			if expr != nil {
				ctx.Args[i] = *expr
			}
		} else {
			ctx.Args[i] = arg
		}
	}

	expr, err := f(ctx)
	if err != nil {
		return nil, err
	}

	return expr, nil
}
