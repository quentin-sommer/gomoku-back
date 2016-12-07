package ia

import (
  "./../protocol"
  "./../referee"
  "fmt"
  "os"
  "log"
  "runtime/pprof"
)

const (
  TWO_ALIGN = 1
  THREE_ALIGN = 50
  FOUR_ALIGN = 100
  // Compute : base + pawn taken
  BASE_PAWN_TAKEN = 4
  // Most important, wins over the rest every time
  FIVE_ALIGN = 5000
  MAX_INIT = -42000
  MIN_INIT = 42000
)

var mapCopies uintptr = 0
var smallestIndex int
var highestIndex int

type minMaxStruct struct {
  M      []protocol.MapData
  Player int8
  Depth  int
  End    bool
  Idx    int
}

func getOtherPlayer(player int8) int8 {
  if (player == 0) {
    return 1
  }
  return 0
}

/**
 * Plays a pawn for player at index idx if possible, otherwise returns false
 * @param  m      map
 * @param  idx    index to play
 * @param  player player to play for
 * @return bool   true if played
 */
func playIdx(m []protocol.MapData, idx int, player int8) bool {
  cell := m[idx]
  if (cell.Empty && cell.Playable) {
    m[idx].Empty = false
    m[idx].Playable = false
    m[idx].Player = player

    return true
  }
  return false
}

func eval(data *minMaxStruct) int {
  ret := 0
  ret += CountSequences(data, 2)
  ret += CountSequences(data, 3)
  ret += CountSequences(data, 4)
  ret += CountSequences(data, 5)

  return ret
/*  one := (TWO_ALIGN * CountSequences(data.M, data.Player, 2))
  two := (THREE_ALIGN * CountSequences(data.M, data.Player, 3))
  three := (FOUR_ALIGN * CountSequences(data.M, data.Player, 4))
  four := 2 * (FIVE_ALIGN * CountSequences(data.M, data.Player, 5))
  if four >= 1 {
    fmt.Println("four ", data.Player, four)
  }
  one2 := (TWO_ALIGN * CountSequences(data.M, getOtherPlayer(data.Player), 2))
  two2 := (THREE_ALIGN * CountSequences(data.M, getOtherPlayer(data.Player), 3))
  three2 := (FOUR_ALIGN * CountSequences(data.M, getOtherPlayer(data.Player), 4))
  four2 := (FIVE_ALIGN * CountSequences(data.M, getOtherPlayer(data.Player), 5))
  if four2 >= 1 {
    fmt.Println("four ",getOtherPlayer(data.Player), four2)
  }
  return ((one + two + three + four) - (one2 + two2 + three2 + four2))*/
}

func max(data *minMaxStruct) int {
  if (data.Depth == 0 || data.End) {
    //data.Player = getOtherPlayer(data.Player)
    return eval(data)
  }
  max := MAX_INIT

  mapcp := make([]protocol.MapData, len(data.M))
  for i := smallestIndex; i < highestIndex; i++ {
    copy(mapcp, data.M)
    mapCopies += 1
    if playIdx(mapcp, i, data.Player) {
      _, end, valid := referee.Exec(mapcp, i)
      if (valid) {
        tmp := min(&minMaxStruct{mapcp, getOtherPlayer(data.Player), data.Depth - 1, end, i})
        if (tmp > max) {
          max = tmp
        }
      }
    }
  }
  return max
}

func min(data *minMaxStruct) int {
  if (data.Depth == 0 || data.End) {
    data.Player = getOtherPlayer(data.Player)
    return eval(data)
  }
  min := MIN_INIT

  mapcp := make([]protocol.MapData, len(data.M))
  for i := smallestIndex; i < highestIndex; i++ {
    copy(mapcp, data.M)
    mapCopies += 1
    if playIdx(mapcp, i, data.Player) {
      _, end, valid := referee.Exec(mapcp, i)
      if (valid) {
        tmp := max(&minMaxStruct{mapcp, getOtherPlayer(data.Player), data.Depth - 1, end, i})
        if (tmp < min) {
        //  fmt.Println("min: ", tmp)
          min = tmp
        }
      }
    }
  }
  return min
}
func initSmallMax(m []protocol.MapData) {
  smallestIndex = -1
  for i := 0; i < protocol.MAP_SIZE; i++ {
    if (!m[i].Empty) {
      if (smallestIndex == -1) {
        smallestIndex = i
      }
      highestIndex = i
    }
  }
  if (smallestIndex > (19 * 2)) {
    smallestIndex -= (19 * 2)
  }
  if ((highestIndex + (19 * 2)) < protocol.MAP_SIZE) {
    highestIndex += (19 * 2)
  }
  // fmt.Println("iteration window size", highestIndex - smallestIndex, "(" + strconv.Itoa(smallestIndex) + "->" + strconv.Itoa(highestIndex) + ")")
}

func MinMaxBenchWrapper(m []protocol.MapData, player int8, depth int) (int) {
  f, err := os.Create("gomoku.prof")
  if err != nil {
    log.Fatal(err)
  }
  pprof.StartCPUProfile(f)
  defer pprof.StopCPUProfile()

  return MinMax(m, player, depth)
}

func MinMax(m []protocol.MapData, player int8, depth int) (int) {
  initSmallMax(m)
  max := MAX_INIT
  maxIdx := 0
  mapCopies = 0

  mapcp := make([]protocol.MapData, len(m))
  for i := smallestIndex; i < highestIndex; i++ {
    copy(mapcp, m)
    mapCopies += 1
    if playIdx(mapcp, i, player) {
      _, end, valid := referee.Exec(mapcp, i)
      if (valid) {
        tmp := min(&minMaxStruct{mapcp, player, depth - 1, end, i})
        if (tmp > max) {
          fmt.Println(tmp, i)
          max = tmp
          maxIdx = i
        }
      }
    }
  }
  // fmt.Println("Total map cp", mapCopies)
  return maxIdx
  //  fmt.Println(m)
  //  fmt.Println(map)
}
