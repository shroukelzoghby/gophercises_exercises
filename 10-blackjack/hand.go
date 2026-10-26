package main

import (
	"strings"

	"github.com/shroukelzoghby/gophercises_exercises/deck"
)

// Hand is the set of cards held by a player or the dealer.
type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

// DealerString renders the dealer's hand with the hole card hidden.
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

// MinScore scores the hand counting every Ace as 1.
func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

// Score returns the best score for the hand, counting one Ace as 11
// when that doesn't bust the hand.
func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			return minScore + 10
		}
	}
	return minScore
}

// Soft reports whether the hand's score counts an Ace as 11.
func (h Hand) Soft() bool {
	return h.Score() != h.MinScore()
}

// Bust reports whether the hand's score is over 21.
func (h Hand) Bust() bool {
	return h.Score() > 21
}

// Blackjack reports whether the hand is a natural 21 (an Ace and a ten-value card).
func (h Hand) Blackjack() bool {
	return len(h) == 2 && h.Score() == 21
}

// draw removes the top card from the deck and returns it with the remaining cards.
func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
