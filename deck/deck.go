//go:generate stringer -type=Suit,Rank
package deck

import (
	"fmt"
	"math/rand"
	"slices"
	"time"
)

type Card struct {
	Rank
	Suit
}

type (
	Option func(Deck) Deck
	Deck   []Card
	Suit   uint8
	Rank   uint8
)

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
	Joker
)

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Returns a new deck
// Takes functional options as arguments to modify
// new decks
func New(opts ...Option) *Deck {
	var deck Deck

	for i := Spades; i <= Hearts; i++ {
		for j := Ace; j <= King; j++ {
			deck = append(deck, Card{Rank: j, Suit: i})
		}
	}

	for _, opt := range opts {
		deck = opt(deck)
	}

	return &deck
}

// Shuffles cards in a deck using a random seed
// and positional shifting
func (d Deck) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range d {
		newPos := r.Intn(len(d) - 1)
		d[i], d[newPos] = d[newPos], d[i]
	}
}

// Returns a string containing a card's Rank and Suit
func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank.String(), c.Suit.String())
}

// Prints each card in a deck, line by line
func (d Deck) PrintDeck() {
	for _, card := range d {
		fmt.Println(card.String())
	}
}

// Functional option for New() function
// Adds n number of jokers to a deck
func WithJokers(n int) func(Deck) Deck {
	return func(d Deck) Deck {
		for i := 0; i < n; i++ {
			d = append(d, Card{Rank: 0, Suit: Joker})
		}
		return d
	}
}

// Functional option for New() function
// Filters cards from deck
func WithoutCard(card Card) func(Deck) Deck {
	var newDeck Deck
	return func(d Deck) Deck {
		newDeck = slices.DeleteFunc(d, func(c Card) bool {
			return c.Rank == card.Rank
		})
		return newDeck
	}
}

// Functional option for New() function
// Provides ability to constuct a deck with multiple decks
func WithMultipleDecks(n int) func(Deck) Deck {
	return func(d Deck) Deck {
		var newDeck Deck
		deck := New()
		for i := 0; i < n; i++ {
			newDeck = append(newDeck, *deck...)
		}
		return newDeck
	}
}
