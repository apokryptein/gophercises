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
	Score  int
}

func main() {
	fmt.Println("Welcome to the BlackJack Game")

	d := deck.New(deck.WithMultipleDecks(3))

	players := []Player{
		{
			Name:   "Player 1",
			Hand:   Deal(2, d),
			Dealer: false,
			Score:  0,
		},
		{
			Name:   "PLayer 2",
			Hand:   Deal(2, d),
			Dealer: false,
			Score:  0,
		},
		{
			Name:   "Dealer",
			Hand:   Deal(2, d),
			Dealer: true,
			Score:  0,
		},
	}

	GameInit(players, d)
	fmt.Println("Results: ", players)
}

func GameInit(players []Player, d *deck.Deck) {
	for idx := range players {
		players[idx].TurnRepl(d)
	}
}

func (p *Player) TurnRepl(d *deck.Deck) {
	p.ScoreHand()
	for {
		p.PrintHand()
		fmt.Printf("Current Score: %d\n\n", p.Score)
		// TODO: Add scoring logic here

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
			p.ScoreHand()
		case "S":
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

func (p *Player) ScoreHand() {
	// TODO: add case for Ace
	score := 0
	for _, c := range p.Hand {
		score += int(cardVal(c))
	}
	p.Score = score
}

func cardVal(c deck.Card) int {
	if c.Rank > 10 {
		return 10
	}
	return int(c.Rank)
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
