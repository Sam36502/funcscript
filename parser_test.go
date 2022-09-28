package funcscript

import (
	"fmt"
	"testing"
)

func TestFizzbuzz(t *testing.T) {
	err := Initialise()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	script := `
		_print(_if(false, "if", _if(false, "elseif", "else")));
		_set("i", 1);
		_while(_not(_eq(_get("i"), 100)),
			_set("i", _sum(_get("i"), 1))
		);
	`
	err = Eval(script)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}
