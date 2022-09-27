package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"fuckyou/funcscript"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type CardList struct {
	Spells  SpellList  `json:"spells"`
	Summons SummonList `json:"summons"`
}

type SpellList struct {
	Level1 []Card `json:"1"`
}

type SummonList struct {
	Level1 []Card `json:"1"`
}

type Card struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (cl CardList) RandomCard() Card {
	allCards := append(cl.Spells.Level1, cl.Summons.Level1...)

	cardslen := len(allCards)
	randIndex := rand.Intn(cardslen)
	card := allCards[randIndex]

	return card
}

func main() {

	err := funcscript.Initialise()
	if err != nil {
		fmt.Println(err)
		return
	}

	TestFunction := func(ctx funcscript.Context) (*funcscript.Expression, error) {
		a, err := ctx.GetDouble(0)
		if err != nil {
			return nil, err
		}

		b, err := ctx.GetDouble(1)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Adding %f + %f\n", a, b)
		return funcscript.DoubleExpr(a + b), nil
	}
	TestCommand := `
		_set("i", 0);
		_if(true, _print("'_if' works"));
		_while(_not(_eq(_get("i"), 10)),
			_set("i", _sum(_get("i"), 1)),
			_print("i: ", _get("i"))
		);
	`

	funcscript.AddFunction("add", TestFunction)

	err = funcscript.Eval(TestCommand)
	if err != nil {
		fmt.Println(err)
		return
	}

	return

	// Read file
	data, err := ioutil.ReadFile("cards.json")
	if err != nil {
		fmt.Println("File no exist; fuck you.", err)
		os.Exit(1)
	}

	// Parse json
	cards := CardList{}
	err = json.Unmarshal(data, &cards)
	if err != nil {
		fmt.Println("Parse fail; fuck you.", err)
		os.Exit(1)
	}

	// Init rand
	rand.Seed(time.Now().UnixMicro())
	rand.Int()

	// Select Cards
	fmt.Println("Sam's Deck:")
	for i := 0; i < 4; i++ {
		fmt.Printf("  %s\n", cards.RandomCard().Name)
	}

	fmt.Println("Weeb's Deck:")
	for i := 0; i < 4; i++ {
		fmt.Printf("  %s\n", cards.RandomCard().Name)
	}

	for {
		fmt.Println("\nPress [ENTER] to get a new card")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Printf("  %s\n", cards.RandomCard().Name)
	}
}
