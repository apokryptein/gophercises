package main

import (
	"fmt"

	"github.com/apokryptein/gophercises/deck"
)

func main() {
	deck := deck.New(deck.WithJokers(2))
	// deck.Shuffle()
	fmt.Println(deck.String())
}
