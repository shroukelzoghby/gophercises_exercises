package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shroukelzoghby/gophercises_exercises/deck"
)

func main() {
	cards := deck.New(deck.Deck(3), deck.Shuffle)
	play(cards, os.Stdin, os.Stdout)
}

// play runs a single round of blackjack between one human player and the dealer.
func play(cards []deck.Card, in io.Reader, out io.Writer) {
	var player, dealer Hand
	var card deck.Card

	// Deal two cards each, alternating player then dealer.
	for i := 0; i < 2; i++ {
		for _, hand := range []*Hand{&player, &dealer} {
			card, cards = draw(cards)
			*hand = append(*hand, card)
		}
	}

	// Player's turn.
	scanner := bufio.NewScanner(in)
	for !player.Bust() {
		fmt.Fprintln(out, "Player:", player, "| score:", player.Score())
		fmt.Fprintln(out, "Dealer:", dealer.DealerString())
		fmt.Fprint(out, "What will you do? (h)it, (s)tand: ")
		if !scanner.Scan() {
			break
		}
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if input == "s" || input == "stand" {
			break
		}
		if input != "h" && input != "hit" {
			fmt.Fprintln(out, "Invalid option:", input)
			continue
		}
		card, cards = draw(cards)
		player = append(player, card)
	}

	// Dealer's turn: hit on 16 or less, or a soft 17. Skipped if the player bust.
	if !player.Bust() {
		for dealer.Score() <= 16 || (dealer.Score() == 17 && dealer.Soft()) {
			card, cards = draw(cards)
			dealer = append(dealer, card)
		}
	}

	fmt.Fprintln(out, "==FINAL HANDS==")
	fmt.Fprintln(out, "Player:", player, "| score:", player.Score())
	fmt.Fprintln(out, "Dealer:", dealer, "| score:", dealer.Score())
	fmt.Fprintln(out, result(player, dealer))
}

// result decides the outcome of a finished round from the player's perspective.
func result(player, dealer Hand) string {
	pScore, dScore := player.Score(), dealer.Score()
	switch {
	case player.Bust():
		return "You busted. Dealer wins."
	case dealer.Bust():
		return "Dealer busted. You win!"
	case pScore > dScore:
		return "You win!"
	case dScore > pScore:
		return "You lose."
	default:
		return "Draw."
	}
}
