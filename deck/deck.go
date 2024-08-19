package deck

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Card struct {
	Value string
	Suit  string
}

type Deck []Card

func New() Deck {
	suits := []string{"Spades", "Diamonds", "Clubs", "Hearts"}
	values := []string{"Ace", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}

	var deck Deck

	for _, suit := range suits {
		for _, val := range values {
			deck = append(deck, Card{val, suit})
		}
	}

	return deck
}

func (d Deck) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range d {
		newPos := r.Intn(len(d) - 1)
		d[i], d[newPos] = d[newPos], d[i]
	}
}

func (d Deck) String() string {
	var deck []string
	for _, card := range d {
		c := fmt.Sprintf("%s of %s", card.Value, card.Suit)
		deck = append(deck, c)
	}

	return strings.Join(deck, ", ")
}
