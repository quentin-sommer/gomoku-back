package ia

import (
  "../protocol"
)

type vec2 struct {
  x, y int
}

func CountSequences(m []protocol.MapData, player int8, seq_len int) int {
  total := 0
  for i := 0; i < protocol.MAP_SIZE; i++ {
    if (m[i].Player == player) {
      total += countSequenceInit(m, i, player, seq_len)
    }
  }
  // fmt.Printf("Total sequence of %d length for player %d : %d\n", seq_len, player, total)

  return total
}

func countSequenceInit(myMap []protocol.MapData, pos int, player int8, seq_len int) int {
  var x int = pos % 19
  var y int = pos / 19
  ret := 0

  tmpmap := make([]protocol.MapData, len(myMap))
  copy(tmpmap, myMap)
  ret += checkSequence(myMap, x, y, &vec2{1, 0}, player, seq_len)
  copy(tmpmap, myMap)
  ret += checkSequence(myMap, x, y, &vec2{0, 1}, player, seq_len)
  copy(tmpmap, myMap)
  ret += checkSequence(myMap, x, y, &vec2{-1, 1}, player, seq_len)
  copy(tmpmap, myMap)
  ret += checkSequence(myMap, x, y, &vec2{1, 1}, player, seq_len)

  return ret
}

func vMove(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX >= 1 && iY == 0) || (protocol.IsInMap(myMap, x + 1, y) && myMap[(x + 1) + y * 19].Player == player))
}

func hMove(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX == 0 && iY >= 1) || (protocol.IsInMap(myMap, x, y + 1) && myMap[x + (y + 1) * 19].Player == player))
}

func D1Move(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX >= 1 && iY >= 1) || (protocol.IsInMap(myMap, x + 1, y + 1) && myMap[(x + 1) + (y + 1) * 19].Player == player))
}

func D2Move(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX <= -1 && iY >= 1) || (protocol.IsInMap(myMap, x - 1, y + 1) && myMap[(x - 1) + (y + 1) * 19].Player == player))
}

func checkSequence(myMap []protocol.MapData, x int, y int, vec *vec2, player int8, seq_len int) int {
  var iX, iY, k int = 0, 0, 1
  iX = vec.x
  iY = vec.y
  for ; k < seq_len && protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == player; {

    if (vMove(myMap, iX, iY, x, y, player) &&
        hMove(myMap, iX, iY, x, y, player) &&
        D1Move(myMap, iX, iY, x, y, player) &&
        D2Move(myMap, iX, iY, x, y, player)) {
      myMap[(x + iX) + (y + iY) * 19].Empty = true
      myMap[(x + iX) + (y + iY) * 19].Player = -1
    }

    k += 1
    iX += vec.x
    iY += vec.y
  }
  if (k >= seq_len) {
    return 1
  }
  return 0
}
