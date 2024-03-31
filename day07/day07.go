package day07

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
)

type cards []int

type hand struct {
	cards cards
	bid   int
}

func parseCards(s string, partb bool) cards {
	rv := make([]int, 0, 5)
	for _, v := range s {
		switch v {
		case 'A':
			rv = append(rv, 14)
		case 'K':
			rv = append(rv, 13)
		case 'Q':
			rv = append(rv, 12)
		case 'J':
			if partb {
				rv = append(rv, 0)
			} else {
				rv = append(rv, 11)
			}
		case 'T':
			rv = append(rv, 10)
		default:
			n := int(v - '0')
			rv = append(rv, n)
		}
	}
	return rv
}

func parseHands(f *os.File, partb bool) []hand {
	var rv []hand

	lines := bufio.NewScanner(f)
	lines.Split(bufio.ScanLines)
	reHand := regexp.MustCompile(`([123456789TJQKA]{5})\s+(\d+)`)
	for lines.Scan() {
		m := reHand.FindStringSubmatch(lines.Text())
		if m == nil {
			log.Fatal("Bad input")
		}
		bid, _ := strconv.Atoi(m[2])
		hand := hand{cards: parseCards(m[1], partb), bid: bid}
		rv = append(rv, hand)
	}

	return rv
}

func getHandPower(hand hand, partb bool) int {
	cards := make([]int, 15)
	for _, cardValue := range hand.cards {
		cards[cardValue]++
	}
	nJokers := cards[0]
	slices.Sort(cards)
	slices.Reverse(cards)
	rv := 0
	switch cards[0] {
	case 5:
		rv = 6
	case 4:
		// 4 of a kind
		rv = 5
		if partb && nJokers == 1 || nJokers == 4 {
			rv++
		}
	case 3:
		switch cards[1] {
		case 2:
			// Full house
			rv = 4
			if partb && (nJokers == 3 || nJokers == 2) {
				// 3 or 2 jokers
				rv = 6
			}
		case 1:
			// Trio
			rv = 3
			if partb && (nJokers == 1 || nJokers == 3) {
				// 4 of a kind
				rv = 5
			}
		}
	case 2:
		switch cards[1] {
		case 2:
			// Two pairs
			rv = 2
			if partb {
				if nJokers == 1 {
					// Full house
					rv = 4
				} else if nJokers == 2 {
					// 4 of a kind
					rv = 5
				}
			}
		case 1:
			// Pair
			rv = 1
			if partb {
				if nJokers == 1 || nJokers == 2 {
					// trio
					rv = 3
				}
			}
		}
	case 1:
		rv = 0
		if partb && nJokers == 1 {
			rv++
		}
	}
	return rv
}

func compareHandPower(lhs, rhs hand, partb bool) int {
	lPower := getHandPower(lhs, partb)
	rPower := getHandPower(rhs, partb)
	if lPower != rPower {
		return lPower - rPower
	}
	for i := range lhs.cards {
		if rv := lhs.cards[i] - rhs.cards[i]; rv != 0 {
			return rv
		}
	}
	log.Fatal("hands equal")
	return 0
}

func compareHandPowerA(lhs, rhs hand) int {
	return compareHandPower(lhs, rhs, false)
}
func compareHandPowerB(lhs, rhs hand) int {
	return compareHandPower(lhs, rhs, true)
}

func Day07a(f *os.File) int {
	rv := 0
	hands := parseHands(f, false)
	slices.SortFunc(hands, compareHandPowerA)
	for i, v := range hands {
		rv += (i + 1) * v.bid
	}
	return rv
}
func Day07b(f *os.File) int {
	rv := 0
	hands := parseHands(f, true)
	slices.SortFunc(hands, compareHandPowerB)
	for i, v := range hands {
		rv += (i + 1) * v.bid
	}
	return rv
}
