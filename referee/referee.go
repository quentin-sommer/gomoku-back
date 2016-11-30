package referee

import (
//"log"
)

import "./../protocol"

const SEAST = 20
const SOUTH = 19
const SWEST = 18
const EAST = 1
const WEST = -1
const NORTH = -19
const NEAST = -18
const NWEST = -20

var Dirtab = [8]int{NORTH, SOUTH, NEAST, SWEST, EAST, WEST, SEAST, NWEST}

func Exec(myMap []protocol.MapData, pos int) ([]protocol.MapData, int, bool, bool) {
	team := myMap[pos].Player
	myMap, capturedPawns := CheckPair(myMap, pos, team)

	ok := CheckDoubleThree(myMap, team)
	if ok == false && capturedPawns > 0 {
		myMap[pos].LegitDoubleThree = true
	} else if ok == false {
		return myMap, 0, false, false
	}

	end := CheckEnd(myMap, pos, team)
	return myMap, capturedPawns, end, true
}

// règles bien expliqué http://maximegirou.com/files/projets/b1/gomoku.pdf

func GetIndexCasePlayed(oldMap []protocol.MapData, newMap []protocol.MapData) int {
  var i int = 0
  for ; i < len(oldMap); i++ {
    if oldMap[i] != newMap[i] {
      return i
    }
  }
  return -1
}

/*
	Function : checkLine
	Parameters :	myMap -> the boardgame with all the pawns
								x -> origin of check on X
								y -> origin of check on Y
								addX -> X of Vector2D
								addY -> Y of Vector2D
								team -> team of the origin pawn to check
	Return : bool -> 5 pawns on a line
	Description:
	Check if there is 5 succent pawns from (x, y) on vector (addX, addY)
*/

func checkLine(myMap []protocol.MapData, x int, y int, addX int, addY int, team int) bool {
  var iX, iY, k int = 0, 0, 1
  iX = addX
  iY = addY
  for ; k < 5 && protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == team; {
    k += 1
    iX += addX
    iY += addY
  }
  iX = -addX
  iY = -addY
  for ; k < 5 && protocol.IsInMap(myMap, x + iX, y + iY) && myMap[(x + iX) + (y + iY) * 19].Player == team; {
    k += 1
    iX -= addX
    iY -= addY
  }
  if (k >= 5) {
    return true
  } else {
    return false
  }
}

func CheckEnd(myMap []protocol.MapData, pos int, team int) bool {
  var x int = pos % 19
  var y int = pos / 19
  if checkLine(myMap, x, y, 1, 0, team) ||
      checkLine(myMap, x, y, 0, 1, team) ||
      checkLine(myMap, x, y, 1, -1, team) ||
      checkLine(myMap, x, y, 1, 1, team) {
    return true
  }
  return false
}

func checkOnPattern(myMap []protocol.MapData, x int, y int, addX int, addY int, team int) bool {
  if (protocol.IsInMap(myMap, x + (addX * 1), y + (addY * 1)) && myMap[x + (addX * 1) + (y + (addY * 1)) * 19].Player == team &&
      protocol.IsInMap(myMap, x + (addX * 2), y + (addY * 2)) && myMap[x + (addX * 2) + (y + (addY * 2)) * 19].Player == team) {
    return true
  }
  if (protocol.IsInMap(myMap, x + (addX * 2), y + (addY * 2)) && myMap[x + (addX * 2) + (y + (addY * 2)) * 19].Player == team &&
      protocol.IsInMap(myMap, x + (addX * 3), y + (addY * 3)) && myMap[x + (addX * 3) + (y + (addY * 3)) * 19].Player == team) {
    return true
  }
  if (protocol.IsInMap(myMap, x - (addX * 1), y - (addY * 1)) && myMap[x - (addX * 1) + (y - (addY * 1)) * 19].Player == team &&
      protocol.IsInMap(myMap, x + (addX * 2), y + (addY * 2)) && myMap[x + (addX * 2) + (y + (addY * 2)) * 19].Player == team) {
    return true
  }
  if (protocol.IsInMap(myMap, x - (addX * 1), y - (addY * 1)) && myMap[x - (addX * 1) + (y - (addY * 1)) * 19].Player == team &&
      protocol.IsInMap(myMap, x + (addX * 1), y + (addY * 1)) && myMap[x + (addX * 1) + (y + (addY * 1)) * 19].Player == team) {
    return true
  }
  return false
}

func CheckDoubleThreeOnOrientation(myMap []protocol.MapData, x int, y int, team int) bool {
  var nbDoubleThree int = 0
  if checkOnPattern(myMap, x, y, 1, 0, team) || checkOnPattern(myMap, x, y, -1, 0, team) {
    nbDoubleThree += 1
  }
  if checkOnPattern(myMap, x, y, 0, 1, team) || checkOnPattern(myMap, x, y, 0, -1, team) {
    nbDoubleThree += 1
  }
  if checkOnPattern(myMap, x, y, 1, -1, team) || checkOnPattern(myMap, x, y, -1, 1, team) {
    nbDoubleThree += 1
  }
  if checkOnPattern(myMap, x, y, 1, 1, team) || checkOnPattern(myMap, x, y, -1, -1, team) {
    nbDoubleThree += 1
  }
  if nbDoubleThree >= 2 {
    return true
  }
  return false
}

func CheckDoubleThree(myMap []protocol.MapData, team int) bool {
	for y:= 0; y < 19; y++ {
		for x:= 0; x < 19; x++ {
			if myMap[x + y * 19].Player == team {
				if myMap[x + y * 19].LegitDoubleThree == false && CheckDoubleThreeOnOrientation(myMap, x, y, team) {
					return false
				}
			}
		}
	}
	return true
}

// function qui check s'il peut NIQUER une paire et s'il peut tej les deux entre (prendre plusieurs pair d'un coup)
func checkCase(myMap []protocol.MapData, pos int, team int) bool {
  if myMap[pos].Player == (team + 1) % 2 {
    return (true)
  }
  return (false)
}

func CheckPair(myMap []protocol.MapData, pos int, team int) ([]protocol.MapData, int) {
  var emptyData protocol.MapData
  emptyData.Empty = true
  emptyData.Playable = true
  emptyData.Player = -1
  captured := 0

  if (pos - (19 * 3)) >= 0 {
    // NORD
    if checkCase(myMap, pos - (19 * 1), team) && checkCase(myMap, pos - (19 * 2), team) && checkCase(myMap, pos - (19 * 3), (team + 1) % 2) {
      myMap[pos - (19 * 1)] = emptyData
      myMap[pos - (19 * 2)] = emptyData
      captured += 2
    }
  }
  if (pos - (19 * 3) + 3) >= 0 && pos % 19 <= 15 {
    // NORD EST
    if checkCase(myMap, pos - (19 * 1) + 1, team) && checkCase(myMap, pos - (19 * 2) + 2, team) && checkCase(myMap, pos - (19 * 3) + 3, (team + 1) % 2) {
      myMap[pos - (19 * 1) + 1] = emptyData
      myMap[pos - (19 * 2) + 2] = emptyData
      captured += 2
    }
  }
  if (pos + 3) < 19 * 19 && pos % 19 <= 15 {
    // EST
    if checkCase(myMap, pos + 1, team) && checkCase(myMap, pos + 2, team) && checkCase(myMap, pos + 3, (team + 1) % 2) {
      myMap[pos + 1] = emptyData
      myMap[pos + 2] = emptyData
      captured += 2
    }
  }
  if (pos + (19 * 3) + 3) < 19 * 19 && pos % 19 <= 15 {
    // SUD EST
    if checkCase(myMap, pos + (19 * 1) + 1, team) && checkCase(myMap, pos + (19 * 2) + 2, team) && checkCase(myMap, pos + (19 * 3) + 3, (team + 1) % 2) {
      myMap[pos + (19 * 1) + 1] = emptyData
      myMap[pos + (19 * 2) + 2] = emptyData
      captured += 2
    }
  }
  if (pos + (19 * 3)) < 19 * 19 {
    // SUD
    if checkCase(myMap, pos + (19 * 1), team) && checkCase(myMap, pos + (19 * 2), team) && checkCase(myMap, pos + (19 * 3), (team + 1) % 2) {
      myMap[pos + (19 * 1)] = emptyData
      myMap[pos + (19 * 2)] = emptyData
      captured += 2
    }
  }
  if (pos + (19 * 3) - 3) < 19 * 19 && pos % 19 >= 3 {
    // SUD OUEST
    if checkCase(myMap, pos + (19 * 1) - 1, team) && checkCase(myMap, pos + (19 * 2) - 2, team) && checkCase(myMap, pos + (19 * 3) - 3, (team + 1) % 2) {
      myMap[pos + (19 * 1) - 1] = emptyData
      myMap[pos + (19 * 2) - 2] = emptyData
      captured += 2
    }
  }
  if (pos - 3) >= 0 && pos % 19 >= 3 {
    // OUEST
    if checkCase(myMap, pos - 1, team) && checkCase(myMap, pos - 2, team) && checkCase(myMap, pos - 3, (team + 1) % 2) {
      myMap[pos - 1] = emptyData
      myMap[pos - 2] = emptyData
      captured += 2
    }
  }
  if (pos - (19 * 3) - 3) >= 0 && pos % 19 >= 3 {
    // NORD OUEST
    if checkCase(myMap, pos - (19 * 1) - 1, team) && checkCase(myMap, pos - (19 * 2) - 2, team) && checkCase(myMap, pos - (19 * 3) - 3, (team + 1) % 2) {
      myMap[pos - (19 * 1) - 1] = emptyData
      myMap[pos - (19 * 2) - 2] = emptyData
      captured += 2
    }
  }
  return myMap, captured
}
