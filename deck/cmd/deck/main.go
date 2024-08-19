package main

import (
	"fmt"

	"github.com/apokryptein/gophercises/deck"
)

func main() {
	deck := deck.New()
	deck.Shuffle()
	fmt.Println(deck.String())
}
