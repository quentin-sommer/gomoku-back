package ia

import (
  //  "./../referee"
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
)

func Min(m []protocol.MapData, player int, profondeur int) (int) {

  if (profondeur == 0) {
    return (1) // score final
  }

  min_val := 0

  for i := 0; i < 19 * 19 ; i++  {

    if m[i].Playable {
      // Simuler coup
      val := Max(m, player, profondeur - 1)

      if (val < min_val && min_val != 0) {
        min_val = val
      }
      // undo le coup
    }
  }


  return -1
}

func Max(m []protocol.MapData, player int, profondeur int) (int) {
  return -1
}

func MinMax(m []protocol.MapData, player int, profondeur int) (int, int) {
  return -1, -1
}