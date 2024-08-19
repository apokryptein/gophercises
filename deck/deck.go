package deck

type Card struct {
	Value string
	Suit  string
}

func New() []Card {
	suits := []string{"Spades", "Diamonds", "Clubs", "Hearts"}
	values := []string{"A", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	var deck []Card

	for _, suit := range suits {
		for _, val := range values {
			deck = append(deck, Card{val, suit})
		}
	}

	return deck
}
