package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	deck "github.com/apokryptein/gophercises/deck"
)

type Player struct {
	Name   string
	Hand   []deck.Card
	Dealer bool
	Stand  bool
	Score  int
}

func main() {
	fmt.Println("Welcome to the BlackJack Game")

	deck := deck.New(deck.WithMultipleDecks(3))

	players := []Player{
		{
			Name:   "Player 1",
			Dealer: false,
			Hand:   Deal(2, deck),
		},
		{
			Name:   "PLayer 2",
			Dealer: false,
			Hand:   Deal(2, deck),
		},
		{
			Name:   "Dealer",
			Dealer: true,
			Hand:   Deal(2, deck),
		},
	}

	GameInit(players, deck)
	fmt.Println("Results: ", players)
}

func GameInit(players []Player, d *deck.Deck) {
	for _, player := range players {
		player.TurnRepl(d)
	}
}

func (p *Player) TurnRepl(d *deck.Deck) {
	for {
		p.PrintHand()

		fmt.Printf("Would you like to Hit (H) or Stand (S)? ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			fmt.Fprintf(os.Stderr, "blackjack: error reading user input: %v\n", err)
			os.Exit(1)
		}

		choice := scanner.Text()

		switch choice {
		case "H":
			p.Hand = slices.Concat(p.Hand, Deal(1, d))
		case "S":
			p.Score = 0 // TODO: implement scoring logic/function
			return
		default:
			fmt.Printf("Invalid choice: %v\n", choice)
		}
	}
}

func Deal(n int, d *deck.Deck) []deck.Card {
	cards := make([]deck.Card, n)
	copy(cards, *d)
	_ = slices.Delete(*d, 0, n)

	return cards
}

func (p *Player) PrintHand() {
	fmt.Printf("\n==== %s ====\n", p.Name)

	if len(p.Hand) == 0 {
		fmt.Printf("NO CARDS IN HAND\n")
		return
	}

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

	fmt.Printf("==================\n\n")
}
