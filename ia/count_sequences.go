package ia

import (
	"github.com/quentin-sommer/gomoku-back/protocol"
)

type vec2 struct {
	x, y int
}

func CountSequences(data *MinMaxStruct, seqLen int) int {
	total := 0
	total += countSequenceInit(data.M, data.Idx, data.Player, seqLen)
	return total
}

func countSequenceInit(myMap []protocol.MapData, pos int, player int8, seq_len int) int {
	var x int = pos % 19
	var y int = pos / 19
	ret := 0

	ret += checkSequence(myMap, x, y, &vec2{1, 0}, player, seq_len)
	ret += checkSequence(myMap, x, y, &vec2{0, 1}, player, seq_len)
	ret += checkSequence(myMap, x, y, &vec2{-1, 1}, player, seq_len)
	ret += checkSequence(myMap, x, y, &vec2{1, 1}, player, seq_len)

	return ret
}

func checkSequence(myMap []protocol.MapData, x int, y int, vec *vec2, player int8, seq_len int) int {
	var iX, iY, k int = 0, 0, 1
	iX = vec.x
	iY = vec.y
	for k < seq_len && protocol.IsInMap(myMap, x+iX, y+iY) && myMap[(x+iX)+(y+iY)*19].Player == player {

		k += 1
		iX += vec.x
		iY += vec.y
	}
	iX = -vec.x
	iY = -vec.y
	for k < seq_len && protocol.IsInMap(myMap, x+iX, y+iY) && myMap[(x+iX)+(y+iY)*19].Player == player {

		k += 1
		iX -= vec.x
		iY -= vec.y
	}

	if k == seq_len {
		return 1
	}
	return 0
}
