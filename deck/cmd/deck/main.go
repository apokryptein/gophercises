package main

import (
	"fmt"

	"github.com/apokryptein/gophercises/deck"
)

func main() {
	deck := deck.New()
	deck.Shuffle()
	deck.PrintDeck()
	fmt.Println("==================================")
	deck.Sort()
	deck.PrintDeck()
}
