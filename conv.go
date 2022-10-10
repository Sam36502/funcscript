package funcscript

import "fmt"

func IntExpr(val int) *Expression {
	return &Expression{
		IntValue: &val,
	}
}

func DoubleExpr(val float64) *Expression {
	return &Expression{
		DoubleValue: &val,
	}
}

func StringExpr(val string) *Expression {
	return &Expression{
		StringValue: &val,
	}
}

func BoolExpr(val bool) *Expression {
	bs := fmt.Sprint(val)
	return &Expression{
		BoolValueStr: &bs,
	}
}
