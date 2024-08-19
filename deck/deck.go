package deck

import (
	"math/rand"
	"time"
)

type Card struct {
	Value string
	Suit  string
}

type Deck []Card

func New() Deck {
	suits := []string{"Spades", "Diamonds", "Clubs", "Hearts"}
	values := []string{"A", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

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
