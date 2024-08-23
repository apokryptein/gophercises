package main

import "github.com/apokryptein/gophercises/deck"

func main() {
	deck := deck.New(deck.WithJokers(2),
		deck.WithMultipleDecks(2),
		deck.WithoutCard(deck.Card{Rank: 2, Suit: 0}))

	deck.PrintDeck()
}
