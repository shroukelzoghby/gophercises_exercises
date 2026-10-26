package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/shroukelzoghby/gophercises_exercises/deck"
)

func c(r deck.Rank) deck.Card {
	return deck.Card{Suit: deck.Spade, Rank: r}
}

func TestScore(t *testing.T) {
	tests := []struct {
		name      string
		hand      Hand
		score     int
		soft      bool
		blackjack bool
	}{
		{"face cards", Hand{c(deck.King), c(deck.Queen)}, 20, false, false},
		{"blackjack", Hand{c(deck.Ace), c(deck.Jack)}, 21, true, true},
		{"soft 17", Hand{c(deck.Ace), c(deck.Six)}, 17, true, false},
		{"hard 17 with ace", Hand{c(deck.Ace), c(deck.Six), c(deck.Ten)}, 17, false, false},
		{"two aces", Hand{c(deck.Ace), c(deck.Ace)}, 12, true, false},
		{"three-card 21", Hand{c(deck.Seven), c(deck.Seven), c(deck.Seven)}, 21, false, false},
		{"bust", Hand{c(deck.King), c(deck.Queen), c(deck.Two)}, 22, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hand.Score(); got != tt.score {
				t.Errorf("Score() = %d, want %d", got, tt.score)
			}
			if got := tt.hand.Soft(); got != tt.soft {
				t.Errorf("Soft() = %v, want %v", got, tt.soft)
			}
			if got := tt.hand.Blackjack(); got != tt.blackjack {
				t.Errorf("Blackjack() = %v, want %v", got, tt.blackjack)
			}
		})
	}
}

func TestPlay(t *testing.T) {
	// Deal order: player, dealer, player, dealer, then draws.
	cards := []deck.Card{
		c(deck.Ten), c(deck.Ace), c(deck.Five), c(deck.Six), // player 15, dealer soft 17
		c(deck.Four),  // player hits to 19
		c(deck.Ten),   // dealer hits soft 17 -> hard 17
		c(deck.Three), // unused: dealer stands on hard 17
	}
	var out bytes.Buffer
	play(cards, strings.NewReader("h\ns\n"), &out)

	got := out.String()
	for _, want := range []string{
		"Player: Ten of Spades, Five of Spades, Four of Spades | score: 19",
		"Dealer: Ace of Spades, Six of Spades, Ten of Spades | score: 17",
		"You win!",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\n%s", want, got)
		}
	}
}

func TestPlayerBustSkipsDealer(t *testing.T) {
	cards := []deck.Card{
		c(deck.Ten), c(deck.Two), c(deck.Nine), c(deck.Three), // player 19, dealer 5
		c(deck.King), // player hits and busts
	}
	var out bytes.Buffer
	play(cards, strings.NewReader("h\n"), &out)

	got := out.String()
	if !strings.Contains(got, "Dealer: Two of Spades, Three of Spades | score: 5") {
		t.Errorf("dealer should not draw after player busts\n%s", got)
	}
	if !strings.Contains(got, "You busted. Dealer wins.") {
		t.Errorf("expected bust result\n%s", got)
	}
}
