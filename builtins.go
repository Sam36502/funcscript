package funcscript

import (
	"fmt"
	"strings"
)

var g_builtins map[string]HandlerFunction = map[string]HandlerFunction{
	"_print": biPrint,
	"_cat":   biCat,
	"_set":   biSet,
	"_get":   biGet,
	"_sum":   biSum,
	"_dif":   biDif,
	"_mul":   biMul,
	"_div":   biDiv,
	"_mod":   biMod,
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

func biCat(ctx Context) (*Expression, error) {
	sb := strings.Builder{}
	for i := range ctx.Args {
		msg, err := ctx.GetAny(i)
		if err != nil {
			return nil, err
		}
		sb.WriteString(msg.String())
	}
	return StringExpr(sb.String()), nil
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

func biDif(ctx Context) (*Expression, error) {
	intA, floatA, err := ctx.GetNumber(0)
	if err != nil {
		return nil, err
	}

	if intA != nil {
		dif := *intA
		for i := 1; i < len(ctx.Args); i++ {
			n, f, err := ctx.GetNumber(i)
			if err != nil {
				return nil, err
			}
			if n != nil {
				dif -= *n
			}
			if f != nil {
				dif -= int(*f)
			}
		}
		return IntExpr(dif), nil
	}

	if floatA != nil {
		dif := *floatA
		for i := 1; i < len(ctx.Args); i++ {
			n, f, err := ctx.GetNumber(i)
			if err != nil {
				return nil, err
			}
			if n != nil {
				dif -= float64(*n)
			}
			if f != nil {
				dif -= *f
			}
		}
		return DoubleExpr(dif), nil
	}

	return nil, nil
}

func biMul(ctx Context) (*Expression, error) {
	intA, floatA, err := ctx.GetNumber(0)
	if err != nil {
		return nil, err
	}

	if intA != nil {
		prod := *intA
		for i := 1; i < len(ctx.Args); i++ {
			n, f, err := ctx.GetNumber(i)
			if err != nil {
				return nil, err
			}
			if n != nil {
				prod *= *n
			}
			if f != nil {
				prod *= int(*f)
			}
		}
		return IntExpr(prod), nil
	}

	if floatA != nil {
		prod := *floatA
		for i := 1; i < len(ctx.Args); i++ {
			n, f, err := ctx.GetNumber(i)
			if err != nil {
				return nil, err
			}
			if n != nil {
				prod *= float64(*n)
			}
			if f != nil {
				prod *= *f
			}
		}
		return DoubleExpr(prod), nil
	}

	return nil, nil
}

func biDiv(ctx Context) (*Expression, error) {
	intA, floatA, err := ctx.GetNumber(0)
	if err != nil {
		return nil, err
	}

	if intA != nil {
		quot := *intA
		for i := 1; i < len(ctx.Args); i++ {
			n, f, err := ctx.GetNumber(i)
			if err != nil {
				return nil, err
			}
			if n != nil {
				if *n == 0 {
					return nil, fmt.Errorf("func '%s' division by 0 not allowed", ctx.FuncName)
				}
				quot /= *n
			}
			if f != nil {
				if *f == 0 {
					return nil, fmt.Errorf("func '%s' division by 0 not allowed", ctx.FuncName)
				}
				quot /= int(*f)
			}
		}
		return IntExpr(quot), nil
	}

	if floatA != nil {
		quot := *floatA
		for i := 1; i < len(ctx.Args); i++ {
			n, f, err := ctx.GetNumber(i)
			if err != nil {
				return nil, err
			}
			if n != nil {
				if *n == 0 {
					return nil, fmt.Errorf("func '%s' division by 0 not allowed", ctx.FuncName)
				}
				quot /= float64(*n)
			}
			if f != nil {
				if *f == 0 {
					return nil, fmt.Errorf("func '%s' division by 0 not allowed", ctx.FuncName)
				}
				quot /= *f
			}
		}
		return DoubleExpr(quot), nil
	}

	return nil, nil
}

func biMod(ctx Context) (*Expression, error) {
	a, err := ctx.GetInt(0)
	if err != nil {
		return nil, err
	}
	b, err := ctx.GetInt(1)
	if err != nil {
		return nil, err
	}
	return IntExpr(a % b), nil
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
	rightExpr, err := ctx.GetAny(1)
	if err != nil {
		return nil, err
	}
	if cond {
		if rightExpr.CommandValue != nil {
			expr, err := evalCommand(*rightExpr.CommandValue)
			if err != nil {
				return expr, err
			}
		} else {
			return &rightExpr, nil
		}
	} else {

		wrongExpr, err := ctx.GetAny(2)
		if err == nil {
			if wrongExpr.CommandValue != nil {
				expr, err := evalCommand(*wrongExpr.CommandValue)
				if err == nil {
					return expr, err
				}
			} else {
				return &wrongExpr, nil
			}
		}
	}
	return nil, nil
}

func biWhile(ctx Context) (*Expression, error) {
	cmds := []Command{}
	for i := 1; i < len(ctx.Args); i++ {
		cmd, err := ctx.getCommand(i)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}

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

		for _, c := range cmds {
			_, err := evalCommand(c)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}
