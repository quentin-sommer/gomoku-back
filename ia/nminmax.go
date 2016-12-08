package ia

import (
  "./../protocol"
  "./../referee"
  "log"
  "fmt"
)

const (
  INF_NEG = -42000
  INF_POS = 42000
)

type infoCell struct {
  worthPlay [2]bool
}

type infoAi struct {
  infoCells [protocol.MAP_SIZE]infoCell
  capturedPawns [2]int
  pawnsPlayed [2]int
  minDepth int
  actualScore int
}

func initInfo(depth int) (*infoAi) {
  info := &infoAi{}
  for i := 0; i < protocol.MAP_SIZE; i++ {
    info.infoCells[i].worthPlay[0] = true
    info.infoCells[i].worthPlay[1] = true
  }
  info.minDepth = depth
  return info
}

func maxVal(a, b int) (int) {
  if a >= b {
    return a
  }
  return b
}

func minVal(a, b int) (int) {
  if a <= b {
    return a
  }
  return b
}

func playPawn(info *infoAi, m []protocol.MapData, player int8, pos int) (int, bool) {
  if info.infoCells[pos].worthPlay[player] && m[pos].Empty && m[pos].Playable {
    m[pos].Empty = false
    m[pos].Player = player
    m[pos].Playable = false
    capturedPawns, _, ok := referee.Exec(m, pos)
    return capturedPawns, ok
  }
  return 0, false
}

func printConsole(m []protocol.MapData) {
  for i := 0; i < protocol.MAP_SIZE; i++ {
    if i % 19 == 0 {
      fmt.Println("")
    }
    if (m[i].Player == 0) {
      fmt.Print("o")
    } else if (m[i].Player == 1) {
      fmt.Print("x")
    } else {
      fmt.Print(".")
    }
  }
  fmt.Println("")
}

func pawnsAround(m []protocol.MapData, pos int, player int8) (int) {
  nb := 0
  x := pos % 19 - 1
  y := pos / 19 - 1
  for i := 0; i < 3; i++ {
    if protocol.IsInMap(m, x + i, y) && !m[(x + i) + y * 19].Empty && m[(x + i) + y * 19].Player != player {
      nb += 1
    }
  }
  x = pos % 19 - 1
  y = pos / 19 + 1
  for i := 0; i < 3; i++ {
    if protocol.IsInMap(m, x + i, y) && !m[(x + i) + y * 19].Empty && m[(x + i) + y * 19].Player != player {
      nb += 1
    }
  }
  y = pos / 19
  x = pos % 19 - 1
  if protocol.IsInMap(m, x, y) && !m[x + y * 19].Empty && m[x + y * 19].Player != player {
    nb += 1
  }
  x = pos % 19 + 1
  if protocol.IsInMap(m, x, y) && !m[x + y * 19].Empty && m[x + y * 19].Player != player {
    nb += 1
  }
  return nb
}

func calculScore(info *infoAi, m []protocol.MapData, pos int, player int8, capturedPawns int, depth int) (int) {
  score := 0
  score += capturedPawns * 100
  //score += pawnsAround(m, pos, player) * 10
  return score * depth
}

func RecurMinMax(info *infoAi, m []protocol.MapData, depth int, playerToMax int8, player int8) (int, int) {
  if depth < info.minDepth {
    info.minDepth = depth
  }
  if depth <= 0 {
    //log.Println(info.capturedPawns, info.pawnsPlayed)
    //printConsole(m)
    //score := info.capturedPawns[playerToMax] * 50 - info.capturedPawns[(playerToMax + 1) % 2] * -(50)
    /*if info.capturedPawns[0] + info.capturedPawns[1] > 0 {
      printConsole(m)
      //log.Println("That was a good move")
      //log.Println(info.capturedPawns, info.pawnsPlayed)
    }*/
    return info.actualScore, 0
  }
  tmpMap := make([]protocol.MapData, protocol.MAP_SIZE)
  if playerToMax == player {
    bestValue := INF_NEG
    bestPos := -1
    for i := 0; i < protocol.MAP_SIZE; i++ {
      copy(tmpMap, m)
      capturedPawns, ok := playPawn(info, tmpMap, player, i)
      if ok {
        info.capturedPawns[player] += capturedPawns
        info.pawnsPlayed[player] += 1
        score := calculScore(info, tmpMap, i, player, capturedPawns, depth)
        info.actualScore += score
        v, _ := RecurMinMax(info, tmpMap, depth - 1, playerToMax, (player + 1) % 2)
        if v >= bestValue {
          bestValue = v
          bestPos = i
        } else {
          info.infoCells[i].worthPlay[player] = false
        }
        info.capturedPawns[player] -= capturedPawns
        info.pawnsPlayed[player] -= 1
        info.actualScore -= score
      }
    }
    return bestValue, bestPos
  } else {
    bestValue := INF_POS
    bestPos := -1
    //log.Println("Player: ", player)
    for i := 0; i < protocol.MAP_SIZE; i++ {
      copy(tmpMap, m)
      capturedPawns, ok := playPawn(info, tmpMap, player, i)
      if ok {
        info.capturedPawns[player] += capturedPawns
        info.pawnsPlayed[player] += 1
        score := calculScore(info, tmpMap, i, player, capturedPawns, depth)
        info.actualScore -= score
        v, _ := RecurMinMax(info, tmpMap, depth - 1, playerToMax, (player + 1) % 2)
        if v <= bestValue {
          bestValue = v
          bestPos = i
        } else {
          info.infoCells[i].worthPlay[player] = false
        }
        info.capturedPawns[player] -= capturedPawns
        info.pawnsPlayed[player] -= 1
        info.actualScore += score
      }
    }
    return bestValue, bestPos
  }
}

func NMinMax(m []protocol.MapData, player int8, depth int) (int) {
  info := initInfo(depth)
  bestScore, bestPos := RecurMinMax(info, m, depth, player, player)
  log.Println("BestScore/BestPos: ", bestScore, " / ", bestPos)
  log.Println("minDepth: ", info.minDepth)
  return bestPos
}
