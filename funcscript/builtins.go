package funcscript

import (
	"fmt"
)

var g_builtins map[string]HandlerFunction = map[string]HandlerFunction{
	"_print": biPrint,
	"_set":   biSet,
	"_get":   biGet,
	"_sum":   biSum,
	"_eq":    biEq,
	"_not":   biNot,
	// "_if":    biIf,    // Commented built-ins are handled in the command
	// "_while": biWhile, // evaluator but are kept here for documentation
}

//func bi(ctx Context) (*Expression, error) {

func biPrint(ctx Context) (*Expression, error) {
	for i := range ctx.Args {
		msg, err := ctx.GetAny(i)
		if err != nil {
			return nil, err
		}
		fmt.Print(msg.String())
	}
	fmt.Print("\n")
	return nil, nil
}

func biSet(ctx Context) (*Expression, error) {
	key, err := ctx.GetString(0)
	if err != nil {
		return nil, err
	}

	val, err := ctx.GetAny(1)
	if err != nil {
		return nil, err
	}

	SetVar(key, val)

	return nil, err
}

func biGet(ctx Context) (*Expression, error) {
	key, err := ctx.GetString(0)
	if err != nil {
		return nil, err
	}

	val, err := GetVar(key)
	if err != nil {
		return nil, err
	}

	return &val, err
}

func biSum(ctx Context) (*Expression, error) {
	intA, floatA, err := ctx.GetNumber(0)
	if err != nil {
		return nil, err
	}

	if intA != nil {
		sum := *intA
		for i := 1; i < len(ctx.Args); i++ {
			n, f, err := ctx.GetNumber(i)
			if err != nil {
				return nil, err
			}
			if n != nil {
				sum += *n
			}
			if f != nil {
				sum += int(*f)
			}
		}
		return IntExpr(sum), nil
	}

	if floatA != nil {
		sum := *floatA
		for i := 1; i < len(ctx.Args); i++ {
			n, f, err := ctx.GetNumber(i)
			if err != nil {
				return nil, err
			}
			if n != nil {
				sum += float64(*n)
			}
			if f != nil {
				sum += *f
			}
		}
		return DoubleExpr(sum), nil
	}

	return nil, nil
}

func biEq(ctx Context) (*Expression, error) {
	a, err := ctx.GetAny(0)
	if err != nil {
		return nil, err
	}
	b, err := ctx.GetAny(1)
	if err != nil {
		return nil, err
	}
	return BoolExpr(a.String() == b.String()), nil
}

func biNot(ctx Context) (*Expression, error) {
	b, err := ctx.GetBool(0)
	if err != nil {
		return nil, err
	}
	return BoolExpr(!b), nil
}

func biIf(ctx Context) (*Expression, error) {
	cond, err := ctx.GetBool(0)
	if err != nil {
		return nil, err
	}
	rightCmd, err := ctx.getCommand(1)
	if err != nil {
		return nil, err
	}
	if cond {
		expr, err := evalCommand(rightCmd)
		if err != nil {
			return expr, err
		}
	}
	wrongCmd, err := ctx.getCommand(1)
	if err != nil {
		return nil, err
	}
	return evalCommand(wrongCmd)
}

func biWhile(ctx Context) (*Expression, error) {
	cond := true
	for cond {
		condExpr, err := ctx.GetAny(0)
		if err != nil {
			return nil, err
		}

		if condExpr.BoolValue != nil {
			cond = *condExpr.BoolValue == "true"
		}
		if condExpr.CommandValue != nil {
			evCond, err := evalCommand(*condExpr.CommandValue)
			if err != nil {
				return nil, err
			}
			if evCond.BoolValue != nil {
				cond = *evCond.BoolValue == "true"
			} else {
				return nil, errMsg("bool", ctx, 0)
			}
		}

		for i := 1; i < len(ctx.Args); i++ {
			cmd, err := ctx.getCommand(i)
			if err != nil {
				return nil, err
			}
			fmt.Println("  cmd: ", cmd.String())
			_, err = evalCommand(cmd)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}
