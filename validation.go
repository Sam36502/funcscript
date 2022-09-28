package funcscript

import "fmt"

type Context struct {
	FuncName string
	Args     []Expression
}

func (ctx Context) GetInt(index int) (int, error) {
	arg := ctx.Args[index].IntValue
	if arg == nil {
		return 0, errMsg("int", ctx, index)
	}
	return *arg, nil
}

func (ctx Context) GetDouble(index int) (float64, error) {
	arg := ctx.Args[index].DoubleValue
	if arg == nil {
		return 0, errMsg("double", ctx, index)
	}
	return *arg, nil
}

func (ctx Context) GetNumber(index int) (*int, *float64, error) {
	arg := ctx.Args[index]
	if arg.DoubleValue != nil {
		return nil, arg.DoubleValue, nil
	}
	if arg.IntValue != nil {
		return arg.IntValue, nil, nil
	}
	return nil, nil, errMsg("number (int/double)", ctx, index)
}

func (ctx Context) GetString(index int) (string, error) {
	arg := ctx.Args[index].StringValue
	if arg == nil {
		return "", errMsg("string", ctx, index)
	}
	return *arg, nil
}

func (ctx Context) GetBool(index int) (bool, error) {
	arg := ctx.Args[index].BoolValue
	if arg == nil {
		return false, errMsg("bool", ctx, index)
	}
	return *arg == "true", nil
}

func (ctx Context) GetAny(index int) (Expression, error) {
	if index < 0 || index >= len(ctx.Args) {
		return Expression{}, fmt.Errorf("func '%s' [%d]: expected any, none received", ctx.FuncName, index)
	}
	return ctx.Args[index], nil
}

//// util functions ////

func (ctx Context) getCommand(index int) (Command, error) {
	arg := ctx.Args[index].CommandValue
	if arg == nil {
		return Command{}, errMsg("command", ctx, index)
	}
	return *arg, nil
}

func errMsg(expected string, ctx Context, index int) error {
	return fmt.Errorf("func '%s' [%d]: expected %s, %s received", ctx.FuncName, index, expected, recType(ctx.Args[index]))
}

// just finds what the received type and arg is for clearer debug messages
func recType(expr Expression) string {
	if expr.IntValue != nil {
		return fmt.Sprintf("int (%d)", *expr.IntValue)
	}
	if expr.DoubleValue != nil {
		return fmt.Sprintf("double (%f)", *expr.DoubleValue)
	}
	if expr.StringValue != nil {
		return fmt.Sprintf("string (%s)", *expr.StringValue)
	}
	if expr.BoolValue != nil {
		return fmt.Sprintf("bool (%s)", *expr.BoolValue)
	}
	if expr.CommandValue != nil {
		return fmt.Sprintf("command (%s)", expr.CommandValue.String())
	}
	return "nothing"
}
