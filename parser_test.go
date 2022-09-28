package funcscript

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	err := Initialise()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func TestHello(t *testing.T) {
	setup()
	script := `
		_print("Hello, world!");
		_print("second", _print("\nfirst"));
	`
	err := Eval(script)
	assert.NoError(t, err)
}

func TestIf(t *testing.T) {
	setup()
	script := `
		_print("0 0: ", _if(false, "if", _if(false, "elseif", "else")));
		_print("0 1: ", _if(false, "if", _if(true,  "elseif", "else")));
		_print("1 0: ", _if(true,  "if", _if(false, "elseif", "else")));
		_print("1 1: ", _if(true,  "if", _if(true,  "elseif", "else")));
	`
	err := Eval(script)
	assert.NoError(t, err)
}

func TestFizzbuzz(t *testing.T) {
	setup()
	script := `
		_set("i", 1);
		_while(_not(_eq(_get("i"), 100)),
			_if(_eq(_mod(_get("i"), 3), 0),
				_print("fizz"),
				_if(_eq(_mod(_get("i"), 5), 0),
					_print("buzz"),
					_if( _eq(_mod(_get("i"), 15), 0),
						_print("fizzbuzz"),
						_print(_get("i"))
					)
				)
			),
			_set("i", _sum(_get("i"), 1))
		);
	`
	err := Eval(script)
	assert.NoError(t, err)
}

func TestCustomFunc(t *testing.T) {
	setup()

	cFunc := func(ctx Context) (*Expression, error) {
		i, e := ctx.GetInt(0)
		assert.NoError(t, e)
		d, e := ctx.GetDouble(1)
		assert.NoError(t, e)
		s, e := ctx.GetString(2)
		assert.NoError(t, e)
		b, e := ctx.GetBool(3)
		assert.NoError(t, e)
		fmt.Printf("Custom function received (%d, %f, '%s', %t)\n", i, d, s, b)
		return nil, nil
	}
	err := AddFunction("custFunc", cFunc)
	assert.NoError(t, err)
	err = Eval(`
		custFunc(69, 3.14159, "Hello, world!", true);
	`)
	assert.NoError(t, err)
}
