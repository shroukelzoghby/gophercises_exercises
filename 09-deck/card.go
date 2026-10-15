
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker // special case; not a real suit
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

// Rank represents a playing card rank (Ace through King).
type Rank uint8

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

const (
	minRank = Ace
	maxRank = King
)

// Card is a single playing card with a suit and rank.
// Jokers have Suit == Joker; their Rank is unused for game logic.
type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}


func New(opts ...func([]Card) []Card) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	for _, opt := range opts {
		cards = opt(cards)
	}
	return cards
}

// DefaultSort sorts cards using Less (suit then rank order).
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}


func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

// Less is the default comparison used with Sort / DefaultSort.
// Absolute order is: all Spades (A..K), then Diamonds, Clubs, Hearts.
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

// Shuffle returns an option that randomly permutes the deck.
func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	perm := r.Perm(len(cards))
	for i, j := range perm {
		ret[i] = cards[j]
	}
	return ret
}

// Jokers returns an option that appends n joker cards to the deck.
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{
				Suit: Joker,
				Rank: Rank(i),
			})
		}
		return cards
	}
}

// Filter returns an option that removes cards for which f returns true.
func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for _, c := range cards {
			if !f(c) {
				ret = append(ret, c)
			}
		}
		return ret
	}
}

// Deck returns an option that repeats the current deck n times
// (e.g. Deck(3) yields three combined standard decks).
func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, cards...)
		}
		return ret
	}
}
