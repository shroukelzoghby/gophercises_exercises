package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Spade})
	fmt.Println(Card{Rank: Nine, Suit: Diamond})
	fmt.Println(Card{Rank: Jack, Suit: Club})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Two of Spades
	// Nine of Diamonds
	// Jack of Clubs
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != 13*4 {
		t.Errorf("wrong number of cards in a new deck: got %d", len(cards))
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	exp := Card{Rank: Ace, Suit: Spade}
	if cards[0] != exp {
		t.Errorf("expected %v as first card; got %v", exp, cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	exp := Card{Rank: Ace, Suit: Spade}
	if cards[0] != exp {
		t.Errorf("expected %v as first card; got %v", exp, cards[0])
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(3))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}
	if count != 3 {
		t.Errorf("expected 3 jokers; got %d", count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(c Card) bool {
		return c.Rank == Two || c.Rank == Three
	}
	cards := New(Filter(filter))
	for _, c := range cards {
		if c.Rank == Two || c.Rank == Three {
			t.Errorf("expected twos and threes to be filtered out")
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))
	if len(cards) != 13*4*3 {
		t.Errorf("expected %d cards; got %d", 13*4*3, len(cards))
	}
}

func TestShuffle(t *testing.T) {
	orig := New()
	shuffled := New(Shuffle)
	if len(shuffled) != len(orig) {
		t.Errorf("shuffle changed deck size: got %d want %d", len(shuffled), len(orig))
	}
}
