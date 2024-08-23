package main

import (
	"fmt"
	"slices"

	deck "github.com/apokryptein/gophercises/deck"
)

type Player struct {
	Name   string
	Hand   []deck.Card
	Dealer bool
}

func main() {
	fmt.Println("Welcome to the BlackJack Game")

	deck := deck.New()

	p1 := Player{
		Name:   "Player 1",
		Dealer: false,
	}

	p1.Hand = Deal(2, deck)

	p1.PrintHand()
}

func Deal(n int, d *deck.Deck) []deck.Card {
	cards := make([]deck.Card, n)
	copy(cards, *d)
	_ = slices.Delete(*d, 0, n)
	return cards
}

func (p *Player) PrintHand() {
	fmt.Printf("==== %s ====\n", p.Name)

	if p.Dealer {
		fmt.Println("HIDDEN CARD")
		for i := 1; i < len(p.Hand); i++ {
			fmt.Println(p.Hand[i])
		}
		return
	} else {
		for _, card := range p.Hand {
			fmt.Println(card)
		}
	}
}
