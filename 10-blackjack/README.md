# Exercise #10: Blackjack Game


## Exercise details

Using the `deck` package from Exercise #9, build a simple command-line blackjack game between one human player and the dealer.

Rules implemented:

1. **Dealing** - the player and dealer are each dealt 2 cards, alternating player then dealer. The dealer's second card stays face down until the end of the round.
2. **Player's turn** - the player chooses to `(h)it` (take another card) or `(s)tand` (end their turn). The turn also ends if the player busts.
3. **Dealer's turn** - the dealer hits on 16 or less, or on a soft 17, and stands otherwise. If the player busted, the dealer doesn't play.
4. **Determining the winner** - the highest score without going over 21 wins. Cards 2-10 are worth face value, J/Q/K are worth 10, and Aces are worth 1 or 11 (whichever is best without busting). A "soft 17" is a 17 in which an Ace is counted as 11.

## Usage

```
go run .
go test ./...
```

The game uses three shuffled decks combined into one shoe.
