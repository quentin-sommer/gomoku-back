package ia

import (
  "../protocol"
)

type vec2 struct {
  x, y int
}

func CountSequences(data *MinMaxStruct, seqLen int) int {
  total := 0
  /*mapcp := make([]protocol.MapData, len(data.M))
  copy(mapcp, data.M)*/
  total += countSequenceInit(data.M, data.Idx, data.Player, seqLen)
  //fmt.Printf("Total sequence of %d length for player %d : %d\n", seqLen, data.Player, total)
  return total
}
/*
func vMove(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX >= 1 && iY == 0) || (protocol.IsInMap(myMap, x + iX, y) && myMap[(x + iX) + y * 19].Player == player))
}

func hMove(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX == 0 && iY >= 1) || (protocol.IsInMap(myMap, x, y + iY) && myMap[x + (y + iY) * 19].Player == player))
}

func d1Move(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX >= 1 && iY >= 1) || (protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == player))
}

func d2Move(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX <= -1 && iY >= 1) || (protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == player))
}

func d1MoveInv(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX <= -1 && iY <= -1) || (protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == player))
}

func d2MoveInv(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX >= 1 && iY <= -1) || (protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == player))
}

func vMoveInv(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX <= -1 && iY == 0) || (protocol.IsInMap(myMap, x + iX, y) && myMap[(x + iX) + y * 19].Player == player))
}

func hMoveInv(myMap []protocol.MapData, iX int, iY int, x int, y int, player int8) bool {
  return ((iX == 0 && iY <= -1) || (protocol.IsInMap(myMap, x, y + iY) && myMap[x + (y + iY) * 19].Player == player))
}
*/

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
  //fmt.Println("x ", x, " y ", y, " iX ", iX, " iY ", iY)
  for ; k < seq_len && protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == player; {

 /*   if (vMove(myMap, iX, iY, x, y, player) &&
        hMove(myMap, iX, iY, x, y, player) &&
        d1Move(myMap, iX, iY, x, y, player) &&
        d2Move(myMap, iX, iY, x, y, player)) {
      myMap[(x + iX) + (y + iY) * 19].Empty = true
      myMap[(x + iX) + (y + iY) * 19].Player = -1
    }*/

    k += 1
   // fmt.Println("Plus 1 dans K dans la premiere boucle", player, x + iX, y + iY)
    iX += vec.x
    iY += vec.y
  }
  iX = -vec.x
  iY = -vec.y
  for ; k < seq_len && protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == player; {

    /*if (vMoveInv(myMap, iX, iY, x, y, player) &&
        hMoveInv(myMap, iX, iY, x, y, player) &&
        d1MoveInv(myMap, iX, iY, x, y, player) &&
        d2MoveInv(myMap, iX, iY, x, y, player)) {
      myMap[(x + iX) + (y + iY) * 19].Empty = true
      myMap[(x + iX) + (y + iY) * 19].Player = -1
    }*/

    k += 1
    //fmt.Println("Plus 1 dans K dans la deuxieme boucle")
    iX -= vec.x
    iY -= vec.y
  }

//  fmt.Println("K value ", k)


  if (k == seq_len) {
    return 1
  }
  return 0
}
