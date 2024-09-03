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

// TODO: Create game state

func main() {
	fmt.Println("Welcome to the BlackJack Game")

	d := deck.New(deck.WithMultipleDecks(3), deck.Shuffle)

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

	dScore := players[2].Score

	fmt.Println("======== RESULTS ========")

	for _, player := range players[0:2] {
		pScore := player.Score

		switch {
		case dScore > 21:
			fmt.Printf("Dealer busted.\n\n")
			return
		case pScore == 21:
			fmt.Printf("%s: BLACKJACK. You win!\n\n", player.Name)
		case pScore > 21:
			fmt.Printf("%s: You busted.\n\n", player.Name)
		case pScore > dScore:
			fmt.Printf("Deal Score: %d\n%s Score: %d\n", dScore, player.Name, pScore)
			fmt.Printf("You win!\n\n")
		case dScore > pScore:
			fmt.Printf("Deal Score: %d\n%s Score: %d\n", dScore, player.Name, pScore)
			fmt.Printf("You lose.\n\n")
		case pScore == dScore:
			fmt.Printf("Deal Score: %d\n%s Score: %d\n", dScore, player.Name, pScore)
			fmt.Printf("Draw.\n\n")
		}
	}
}

func GameInit(players []Player, d *deck.Deck) {
	for idx := range players {
		if players[idx].Dealer {
			players[idx].ScoreHand()
			// TODO: account for soft 17
			for players[idx].Score <= 16 {
				players[idx].Hand = slices.Concat(players[idx].Hand, Deal(1, d))
				players[idx].ScoreHand()
			}
			continue
		}
		players[idx].TurnRepl(d)
	}
}

func (p *Player) TurnRepl(d *deck.Deck) {
	p.ScoreHand()
	for {
		p.PrintHand()
		fmt.Printf("Current Score: %d\n\n", p.Score)
		if p.Score > 21 {
			fmt.Println("You busted.")
			return
		}
		// TODO: show dealer's hand here to aid in decision making
		fmt.Printf("Would you like to Hit (h) or Stand (s)? ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			fmt.Fprintf(os.Stderr, "blackjack: error reading user input: %v\n", err)
			os.Exit(1)
		}

		choice := scanner.Text()

		switch choice {
		case "h":
			p.Hand = slices.Concat(p.Hand, Deal(1, d))
			p.ScoreHand()
		case "s":
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

// Calculates score of hand
func (p *Player) ScoreHand() {
	minScore := MinScore(p)

	// if min score is greater than one
	// we cannot count Aces as having a value of 11
	if minScore > 11 {
		p.Score = minScore
		return
	}

	// Check for Aces
	// if minScore is less than 11 and an Ace
	// is present, add 10 to score
	for _, card := range p.Hand {
		if card.Rank == deck.Ace {
			p.Score = minScore + 10
			return
		}
	}
	p.Score = minScore
}

func cardVal(c deck.Card) int {
	if c.Rank > 10 {
		return 10
	}
	return int(c.Rank)
}

// Calculates the score assuming all Aces
// have a value of 1
func MinScore(p *Player) int {
	score := 0
	for _, c := range p.Hand {
		score += int(cardVal(c))
	}
	return score
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
