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

func (d Deck) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range d {
		newPos := r.Intn(len(d) - 1)
		d[i], d[newPos] = d[newPos], d[i]
	}
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank.String(), c.Suit.String())
}

func (d Deck) PrintDeck() {
	for _, card := range d {
		fmt.Println(card.String())
	}
}

func WithJokers(n int) func(Deck) Deck {
	return func(d Deck) Deck {
		for i := 0; i < n; i++ {
			d = append(d, Card{Rank: 0, Suit: Joker})
		}
		return d
	}
}

func WithoutCard(card Card) func(Deck) Deck {
	var newDeck Deck
	return func(d Deck) Deck {
		newDeck = slices.DeleteFunc(d, func(c Card) bool {
			return c.Rank == card.Rank
		})
		return newDeck
	}
}

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
