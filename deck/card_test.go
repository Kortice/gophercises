package deck

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Spade})
	fmt.Println(Card{Rank: Jack, Suit: Club})
	fmt.Println(Card{Rank: King, Suit: Diamond})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Two of Spades
	// Jack of Clubs
	// King of Diamonds
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	// 13 ranks * 4 suits
	want := 13 * 4
	if len(cards) != want {
		t.Errorf("got cards length: %d, but want %d", len(cards), want)
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	want := Card{Rank: Ace, Suit: Spade}

	if cards[0] != want {
		t.Errorf("got %+v, but want %+v", cards[0], want)
	}
}

func TestSort(t *testing.T) {
	// cards := New(Sort(func(cards []Card) func(i int, j int) bool {
	// 	return func(i, j int) bool {
	// 		return absRank(cards[i]) < absRank(cards[j])
	// 	}
	// }))
	cards := New(Sort(Less))
	want := Card{Rank: Ace, Suit: Spade}

	if cards[0] != want {
		t.Errorf("got %+v, but want %+v", cards[0], want)
	}
}

func TestShuffle(t *testing.T) {
	orig := New()

	shuffleRand = rand.New(rand.NewSource(0))
	perm := shuffleRand.Perm(len(orig))

	shuffleRand = rand.New(rand.NewSource(0))
	cards := New(Shuffle)

	for i, j := range perm {
		if cards[i] != orig[j] {
			t.Error("shuffle func got wrong!")
		}
	}
}

func TestJokers(t *testing.T) {
	jokerNum := 3
	cards := New(Jokers(jokerNum))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}
	if count != jokerNum {
		t.Errorf("got %d, but want %d", count, jokerNum)
	}
}

func TestFilter(t *testing.T) {
	filter := func(c Card) bool {
		return c.Rank == Two || c.Rank == Three
	}
	cards := New(Filter(filter))
	for _, card := range cards {
		if card.Rank == Two || card.Rank == Three {
			t.Error("filter failed!")
		}
	}
}
