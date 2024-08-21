package main

import (
	"fmt"

	"github.com/apokryptein/gophercises/deck"
)

func main() {
	deck := deck.New(deck.WithoutCard(deck.Card{Rank: 2, Suit: 0}),
		deck.WithoutCard(deck.Card{Rank: 3, Suit: 0}))
	fmt.Println(deck.String())
}
