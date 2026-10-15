
package deck

import "strconv"

func _() {
	
	var x [1]struct{}
	_ = x[Spade-0]
	_ = x[Diamond-1]
	_ = x[Club-2]
	_ = x[Heart-3]
	_ = x[Joker-4]
}

const _Suit_name = "SpadeDiamondClubHeartJoker"

var _Suit_index = [...]uint8{0, 5, 12, 16, 21, 26}

func (i Suit) String() string {
	idx := int(i) - 0
	if i < 0 || idx >= len(_Suit_index)-1 {
		return "Suit(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Suit_name[_Suit_index[idx]:_Suit_index[idx+1]]
}
func _() {
	
	var x [1]struct{}
	_ = x[Ace-1]
	_ = x[Two-2]
	_ = x[Three-3]
	_ = x[Four-4]
	_ = x[Five-5]
	_ = x[Six-6]
	_ = x[Seven-7]
	_ = x[Eight-8]
	_ = x[Nine-9]
	_ = x[Ten-10]
	_ = x[Jack-11]
	_ = x[Queen-12]
	_ = x[King-13]
}

const _Rank_name = "AceTwoThreeFourFiveSixSevenEightNineTenJackQueenKing"

var _Rank_index = [...]uint8{0, 3, 6, 11, 15, 19, 22, 27, 32, 36, 39, 43, 48, 52}

func (i Rank) String() string {
	idx := int(i) - 1
	if i < 1 || idx >= len(_Rank_index)-1 {
		return "Rank(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Rank_name[_Rank_index[idx]:_Rank_index[idx+1]]
}
