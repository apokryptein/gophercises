package main

import (
	"fmt"
	"strings"

	"github.com/apokryptein/gophercises/deck"
)

func main() {
	deck := deck.New(deck.WithMultipleDecks(2),
		deck.WithJokers(2),
		deck.WithoutCard(deck.Card{Rank: 2, Suit: 0}))

	fmt.Println(strings.ReplaceAll(deck.String(), ", ", "\n"))
}
