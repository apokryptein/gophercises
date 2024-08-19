package main

import (
	"fmt"

	"github.com/apokryptein/gophercises/deck"
)

func main() {
	deck := deck.New()
	fmt.Println(deck)

	deck.Shuffle()
	fmt.Println(deck)
}
