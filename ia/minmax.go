package ia

import (
  "./../referee"
  "./../protocol"
)

/*
 * We start off by guiding ai to build up sequences, and then emphasis on taking pawns
 */
const (
  TWO_ALIGN = 1
  THREE_ALIGN = 2
  FOUR_ALIGN = 3
  // Compute : base + pawn taken
  BASE_PAWN_TAKEN = 4
  // Most important, wins over the rest every time
  FIVE_ALIGN = 500
  MAP_SIZE = 19 * 19
  NON_INIT = -42
)

func eval(m []protocol.MapData, player int) (int) {
  val := 0



  return val
}

func Min(m []protocol.MapData, player int, depth int) (int) {

  if (depth == 0) {
    return (1) // return value of the move
  }

  tmpmap := make([]protocol.MapData, len(m))
  copy(tmpmap, m)
  min_val := NON_INIT
  ok := false

  for i := 0; i < MAP_SIZE; i++ {
    // Simuler coup
    tmpmap, _, _, ok = referee.Exec(tmpmap, i)

    if (ok) {

      val := Max(tmpmap, player, depth - 1)

      if (val < min_val || min_val == NON_INIT) {
        min_val = val
      }
    }
  }
  return min_val
}

func Max(m []protocol.MapData, player int, depth int) (int) {

  if (depth == 0) {
    return (1) // return value of the move
  }

  tmpmap := make([]protocol.MapData, len(m))
  copy(tmpmap, m)
  max_val := NON_INIT
  ok := false

  for i := 0; i < MAP_SIZE; i++ {
    // Simuler coup
    tmpmap, _, _, ok = referee.Exec(tmpmap, i)

    if (ok) {
      val := Min(tmpmap, player, depth - 1)

      if (val > max_val || max_val == NON_INIT) {
        max_val = val
      }
    }
  }
  return max_val
}

func MinMax(m []protocol.MapData, player int, profondeur int) (int, int) {
  return -1, -1
}
