# funcscript
*extremely basic script language for easily extensible scripts*

This library aims to generically cover the cases where you might want
to be able to store functions that can easily be expanded in the future.

I decided to write this as a possible way to handle the effects of
playing cards in a card-game engine where I wasn't sure how the effects
might develop as the game developed. Initially, we imagined a system of
storing function names and arguments, which I decided to expand into a
very basic scripting language. This would allow new functions to be added
or changed easily, without having to completely rewrite the data structures.

Here's an example of how you might handle a card that, once played, will
attack a random target for 10 damage every turn for 5 turns:

	_set("trgt", randomTarget());
	repeatOnTurn(5, dmg(_get("trgt"), 10));

## Usage
Here's a basic example of a possible use-case, assuming `DamageTarget()` is defined somewhere:
```go

package main

func main() {
	err := funcscript.Initialise()
	if err != nil {
		fmt.Println("Failed to initialise:", err)
	}
	
	err = AddFunction("dmg", func(ctx Context) (*Expression, error) {
		targetname, err := ctx.GetString(0)
		if err != nil {
			return nil, err
		}
		dmgAmount, err := ctx.GetInt(1)
		if err != nil {
			return nil, err
		}

		DamageTarget(targetname, dmgAmount)
		return nil, nil
	})
	if err != nil {
		fmt.Println("Failed to add dmg function:", err)
	}
	
	err = Eval(`
		_set("target", "steve");
		dmg(_get("target"), 10);
	`)
	if err != nil {
		fmt.Println("Failed to evaluate script:", err)
	}
}

```

## FuncScript syntax
It should be fairly self-explanatory if you're familiar with C-like syntaxes.
Essentially, it's a subset of C which only supports function calls.

Every script consists of a list of commands ended with semicolon (;).
A command then simply consists of a name followed by arguments in parentheses.
(e.g. `cmdName("string arg", 42, 3.14, true, subcmd());`)

Then arguments can be any of: integer, double, string, boolean or another command,
which will be evaluated before passing in its value.

## Built-In Functions
Generally, most of the use of FuncScript is to call user-defined functions, but for
convenience and debugging, some handy functions for general use. To help distinguish
these, they are all prefixed with an underscore (_):

| Function name          | Description                                                                                           |
|------------------------|-------------------------------------------------------------------------------------------------------|
| `_print(Any...)`       | Takes any number of any type and prints them all to stdout with a newline attached.                   |
| `_cat(Any...)`         | Takes any number of any type and returns a string with them all concatenated.                         |
| `_set(String, Any)`    | Takes a key and value type and stores that in a global store.                                         |
| `_get(String)`         | Takes a key and fetches the value from the global store. Errors if the key doesn't exist.             |
| `_sum(Number...)`      | Takes any number of Ints/Doubles and returns their sum. Return type is same as first argument.        |
| `_dif(Number...)`      | Takes any number of Ints/Doubles and returns their difference. Return type is same as first argument. |
| `_mul(Number...)`      | Takes any number of Ints/Doubles and returns their product. Return type is same as first argument.    |
| `_div(Number...)`      | Takes any number of Ints/Doubles and returns their quotient. Return type is same as first argument.   |
| `_mod(Int, Int)`       | Takes two numbers and returns their modulus.                                                          |
| `_eq(Any, Any)`        | Takes two arguments of any type and returns whether they are equal. Very lenient ('5' == 5.0 == 5).   |
| `_not(Bool)`           | Takes a boolean and returns its inverse (same as '!' operator in most languages).                     |
| `_if(Bool, Any, Any?)` | Takes a condition, a true value and a false value. False value is optional.                           |
| `_while(Bool, Cmd...)` | Takes a condition and repeats the following commands as long as the condition is true.                |

## Issues to specify
 1. Type system? (stringly?, roll int and double into 'number'?, etc.), Value type
 2. Basic mathematical expressions (e.g. "a + b - 4"), vars like "a = bla() + 5"
 3. Breadth-first command evaluation
 4. Error-handling & debug info
 5. Customisable syntax (change separators, custom lexer?, etc.)
 6. ~~Returning results from `Eval()` (`_return()` ?)~~
 7. Performance concerns (pre-build parser?)
 8. Define new functions in the script? (_func("cmdname", cmds...))
