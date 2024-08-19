//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Card struct {
	Rank
	Suit
}

type (
	Deck []Card
	Suit uint8
	Rank uint8
)

const (
	Spade Suit = iota
	Diamonds
	Clubs
	Hearts
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

func New() *Deck {
	var deck Deck

	for i := Spade; i <= Hearts; i++ {
		for j := Ace; j <= King; j++ {
			deck = append(deck, Card{j, i})
		}
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

// TODO: maybe implement this function as a Card to String
// in stead of entire deck to string
func (d Deck) String() string {
	var deck []string
	for _, card := range d {
		c := fmt.Sprintf("%s of %s", card.Rank.String(), card.Suit.String())
		deck = append(deck, c)
	}

	return strings.Join(deck, ", ")
}
