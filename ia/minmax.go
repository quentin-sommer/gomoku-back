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
  NON_INIT = -42
)

func eval(m []protocol.MapData, player int, capture int) (int) {
  val := 0

  return val
}

func min(m []protocol.MapData, player int, depth int, capture int) (int) {

  if (depth == 0) {
    return eval(m, player, capture) // return value of the move
  }

  tmpmap := make([]protocol.MapData, len(m))
  copy(tmpmap, m)
  min_val := NON_INIT
  ok := false

  for i := 0; i < protocol.MAP_SIZE; i++ {
    // Simuler coup
    tmpmap, capture, _, ok = referee.Exec(tmpmap, i)

    if (ok) {

      val := max(tmpmap, player, depth - 1, capture)

      if (val < min_val || min_val == NON_INIT) {
        min_val = val
      }
    }
  }
  return min_val
}

func max(m []protocol.MapData, player int, depth int, capture int) (int) {

  if (depth == 0) {
    return eval(m, player, capture) // return value of the move
  }

  tmpmap := make([]protocol.MapData, len(m))
  copy(tmpmap, m)
  max_val := NON_INIT
  ok := false

  for i := 0; i < protocol.MAP_SIZE; i++ {
    // Simuler coup
    tmpmap, capture, _, ok = referee.Exec(tmpmap, i)
    if (ok) {
      val := min(tmpmap, player, depth - 1, capture)
      if (val > max_val || max_val == NON_INIT) {
        max_val = val
      }
    }
  }
  return max_val
}

func MinMax(m []protocol.MapData, player int, depth int) (int, int) {
  ret := max(m, player, depth, 0)
  CountSequences(m, player, 2)
  return ret, -1
}
