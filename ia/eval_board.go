package ia

import (
  "../protocol"
  "fmt"
)

type vec2 struct {
  x, y int
}

func CountSequences(m []protocol.MapData, player int, seq_len int) int {
  tmpmap := make([]protocol.MapData, len(m))
  copy(tmpmap, m)
  total := 0
  for i := 0; i < protocol.MAP_SIZE; i++ {
    if (tmpmap[i].Player == player) {
      total += countSequenceInit(tmpmap, i, player, seq_len)
    }
  }
  fmt.Printf("Total sequence of %d length for player %d : %d\n", seq_len, player, total)

  return total
}

func countSequenceInit(myMap []protocol.MapData, pos int, player int, seq_len int) int {
  var x int = pos % 19
  var y int = pos / 19
  ret := 0

  ret += checkSequence(myMap, x, y, &vec2{1, 0}, player, seq_len)
  ret += checkSequence(myMap, x, y, &vec2{0, 1}, player, seq_len)
  ret += checkSequence(myMap, x, y, &vec2{1, -1}, player, seq_len)
  ret += checkSequence(myMap, x, y, &vec2{1, 1}, player, seq_len)

  return ret
}

func checkSequence(myMap []protocol.MapData, x int, y int, vec *vec2, player int, seq_len int) int {
  var iX, iY, k int = 0, 0, 1
  iX = vec.x
  iY = vec.y
  for ; k < seq_len && protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == player; {
    myMap[(x + iX) + (y + iY) * 19].Empty = true
    myMap[(x + iX) + (y + iY) * 19].Player = -1

    k += 1
    iX += vec.x
    iY += vec.y
  }
  iX = -vec.x
  iY = -vec.y
  for ; k < seq_len && protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == player; {
    myMap[(x + iX) + (y + iY) * 19].Empty = true
    myMap[(x + iX) + (y + iY) * 19].Player = -1

    k += 1
    iX -= vec.x
    iY -= vec.y
  }
  if (k >= seq_len) {
    return 1
  } else {
    return 0
  }
}
