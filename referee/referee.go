package referee

import (
	"./../protocol"
	"fmt"
)

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
	ok := Checkdoublethree(myMap, pos, team)
	if ok == false {
		return myMap, 0, false, ok
	}
	myMap, capturedPawns := CheckPair(myMap, pos, team)
	end := CheckEnd(myMap, pos, team)
	return myMap, capturedPawns, end, ok
}

// règles bien expliqué http://maximegirou.com/files/projets/b1/gomoku.pdf

// function qui check dans un sens choisi (N NE E SE S SW W NW) pour vérifier la fin du jeu
// attention un gars peut casser une ligne de 5 pions avec une paire

func checkLigne(myMap []protocol.MapData, pos int, team int, val int, add int) int {
	if pos < 19*19 && pos >= 0 {
		if (add == -18 && pos % 19 != 18) || (add == 18 && pos % 19 != 0) || (add == -20 && pos % 19 != 0) || (add == 20 && pos % 19 != 18) || add == 1 || add == -1 || add == 19 || add == -19 {
			if myMap[pos].Player == team {
				return checkLigne(myMap, pos + add, team, val + 1, add)
			}
		}
	}
	return val
}

func CheckEnd(myMap []protocol.MapData, pos int, team int) bool {
	nb := 0
	// horizontal
	nb = checkLigne(myMap, pos, team, 0, 1)
	nb += checkLigne(myMap, pos, team, 0, -1)
	if nb-1 == 5 {
		//fmt.Printf("END 5 IN A ROW 1\n")
		return true
	}

	// vertical
	nb = checkLigne(myMap, pos, team, 0, 19)
	nb += checkLigne(myMap, pos, team, 0, -19)
	if nb-1 == 5 {
		//fmt.Printf("END 5 IN A ROW 2\n")
		return true
	}

	// diagonal /
	nb = checkLigne(myMap, pos, team, 0, -18)
	nb += checkLigne(myMap, pos, team, 0, 18)
	if nb-1 == 5 {
		//fmt.Printf("END 5 IN A ROW 3\n")
		return true
	}

	// diagonal \
	nb = checkLigne(myMap, pos, team, 0, -20)
	nb += checkLigne(myMap, pos, team, 0, 20)
	if nb-1 == 5 {
		//fmt.Printf("END 5 IN A ROW 4\n")
		return true
	}
	return false
}

func GetIndexCasePlayed(oldMap []protocol.MapData, newMap []protocol.MapData) int {
	var i int = 0
	for ; i < len(oldMap); i++ {
		if oldMap[i] != newMap[i] {
			return i
		}
	}
	return -1
}

// function qui check la regle "LE DOUBLE-TROIS"

func getIndexWithDir(pos int, team int, dir int) int {
	//check out of map
	if (pos+dir) > (19*19) || (pos+dir) < 0 {
		return -1
	}

	if (dir == SOUTH || dir == NORTH) && ((pos+dir)%19) == (pos%19) {

		return (pos + dir)

	} else if (dir == EAST || dir == WEST) && ((pos+dir)/19) == (pos/19) {

		return (pos + dir)

	} else if (dir == NEAST || dir == NWEST) && ((pos+dir)/19) == ((pos/19)-1) {

		if dir == NEAST && ((pos+dir)%19) == ((pos%19)+1) {
			return (pos + dir)
		} else if dir == NWEST && ((pos+dir)%19) == ((pos%19)-1) {
			return (pos + dir)

		}

	} else if (dir == SEAST || dir == SWEST) && ((pos+dir)/19) == ((pos/19)+1) {

		if dir == SEAST && ((pos+dir)%19) == ((pos%19)+1) {
			return (pos + dir)
		} else if dir == SWEST && ((pos+dir)%19) == ((pos%19)-1) {
			return (pos + dir)
		}

	}
	return -1
}

func getNbPionTeamIndir(myMap []protocol.MapData, pos int, team int, dir int) int {
	var idx = getIndexWithDir(pos, team, dir)

	if idx < 0 || idx >= (19*19) {
		return 0
	}
	if myMap[idx].Player == team {
		return 1 + getNbPionTeamIndir(myMap, idx, team, dir)
	} else if myMap[idx].Player == -1 {
		return 0
	}
	return -4
}

func Checkdoublethree(myMap []protocol.MapData, pos int, team int) bool {
	var nbline = 0
	var nbdiag = 0
	var pass = 0

	for i := 0; i < 8; i++ {
		if nbline != 2 && myMap[Dirtab[i]+pos].Player == team {
			nbline = getNbPionTeamIndir(myMap, pos, team, Dirtab[i])
			if nbline == 2 {
				pass = 1
			}
		}
		if pass == 0 && nbdiag != 2 && myMap[pos+Dirtab[i]].Player == -1 && myMap[pos+(Dirtab[i]*2)].Player == team {
			nbdiag = getNbPionTeamIndir(myMap, pos+(Dirtab[i]), team, Dirtab[i])
		}
		if nbline == 2 && nbdiag == 2 {
			return false
		}
		pass = 0
	}
	return true
}

// function qui check s'il peut NIQUER une paire et s'il peut tej les deux entre (prendre plusieurs pair d'un coup)
func checkCase(myMap []protocol.MapData, pos int, team int) bool {
	if myMap[pos].Player == (team+1)%2 {
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
		if checkCase(myMap, pos-(19*1), team) && checkCase(myMap, pos-(19*2), team) && checkCase(myMap, pos-(19*3), (team+1)%2) {
			myMap[pos-(19*1)] = emptyData
			myMap[pos-(19*2)] = emptyData
			captured += 2
		}
	}
	if (pos-(19*3)+3) >= 0 && pos%19 <= 15 {
		// NORD EST
		if checkCase(myMap, pos-(19*1)+1, team) && checkCase(myMap, pos-(19*2)+2, team) && checkCase(myMap, pos-(19*3)+3, (team+1)%2) {
			myMap[pos-(19*1)+1] = emptyData
			myMap[pos-(19*2)+2] = emptyData
			captured += 2
		}
	}
	if (pos+3) < 19*19 && pos%19 <= 15 {
		// EST
		if checkCase(myMap, pos+1, team) && checkCase(myMap, pos+2, team) && checkCase(myMap, pos+3, (team+1)%2) {
			myMap[pos+1] = emptyData
			myMap[pos+2] = emptyData
			captured += 2
		}
	}
	if (pos+(19*3)+3) < 19*19 && pos%19 <= 15 {
		// SUD EST
		if checkCase(myMap, pos+(19*1)+1, team) && checkCase(myMap, pos+(19*2)+2, team) && checkCase(myMap, pos+(19*3)+3, (team+1)%2) {
			myMap[pos+(19*1)+1] = emptyData
			myMap[pos+(19*2)+2] = emptyData
			captured += 2
		}
	}
	if (pos + (19 * 3)) < 19*19 {
		// SUD
		if checkCase(myMap, pos+(19*1), team) && checkCase(myMap, pos+(19*2), team) && checkCase(myMap, pos+(19*3), (team+1)%2) {
			myMap[pos+(19*1)] = emptyData
			myMap[pos+(19*2)] = emptyData
			captured += 2
		}
	}
	if (pos+(19*3)-3) < 19*19 && pos%19 >= 2 {
		// SUD OUEST
		if checkCase(myMap, pos+(19*1)-1, team) && checkCase(myMap, pos+(19*2)-2, team) && checkCase(myMap, pos+(19*3)-3, (team+1)%2) {
			myMap[pos+(19*1)-1] = emptyData
			myMap[pos+(19*2)-2] = emptyData
			captured += 2
		}
	}
	if (pos-3) >= 0 && pos%19 >= 2 {
		// OUEST
		if checkCase(myMap, pos-1, team) && checkCase(myMap, pos-2, team) && checkCase(myMap, pos-3, (team+1)%2) {
			myMap[pos-1] = emptyData
			myMap[pos-2] = emptyData
			captured += 2
		}
	}
	if (pos-(19*3)-3) >= 0 && pos%19 >= 2 {
		// NORD OUEST
		if checkCase(myMap, pos-(19*1)-1, team) && checkCase(myMap, pos-(19*2)-2, team) && checkCase(myMap, pos-(19*3)-3, (team+1)%2) {
			myMap[pos-(19*1)-1] = emptyData
			myMap[pos-(19*2)-2] = emptyData
			captured += 2
		}
	}
	return myMap, captured
}
