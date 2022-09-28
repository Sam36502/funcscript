package funcscript

import "fmt"

type Context struct {
	FuncName string
	Args     []Expression
}

func (ctx Context) GetInt(index int) (int, error) {
	arg := ctx.Args[index]
	if arg.IntValue == nil {
		expr, err := ctx.evalCmdArg(index)
		if err != nil {
			return 0, errMsg("int", ctx, index)
		}
		return *expr.IntValue, nil
	}
	return *arg.IntValue, nil
}

func (ctx Context) GetDouble(index int) (float64, error) {
	arg := ctx.Args[index]
	if arg.DoubleValue == nil {
		expr, err := ctx.evalCmdArg(index)
		if err != nil {
			return 0, errMsg("double", ctx, index)
		}
		return *expr.DoubleValue, nil
	}
	return *arg.DoubleValue, nil
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
	arg := ctx.Args[index]
	if arg.StringValue == nil {
		expr, err := ctx.evalCmdArg(index)
		if err != nil {
			return "", errMsg("string", ctx, index)
		}
		return *expr.StringValue, nil
	}
	return *arg.StringValue, nil
}

func (ctx Context) GetBool(index int) (bool, error) {
	arg := ctx.Args[index]
	if arg.BoolValue == nil {
		expr, err := ctx.evalCmdArg(index)
		if err == nil {
			return *expr.BoolValue == "true", nil
		}

		// Try coercing
		if arg.IntValue != nil {
			return *arg.IntValue > 0, nil
		}
		if arg.DoubleValue != nil {
			return *arg.DoubleValue > 0, nil
		}
		if arg.StringValue != nil {
			return *arg.StringValue == "true", nil
		}
		return false, errMsg("bool", ctx, index)
	}
	return *arg.BoolValue == "true", nil
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

// Formats error messages
func errMsg(expected string, ctx Context, index int) error {
	return fmt.Errorf("func '%s' [%d]: expected %s, %s received", ctx.FuncName, index, expected, recType(ctx.Args[index]))
}

// Evaluates commands when other types were expected
func (ctx *Context) evalCmdArg(index int) (*Expression, error) {
	cmd, err := ctx.getCommand(index)
	if err != nil {
		return nil, err
	}
	result, err := evalCommand(cmd)
	if err != nil {
		return nil, err
	}
	return result, nil
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
