package referee

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
	if pos < 19 * 19 && pos >= 0 {
		if (add == -18 && pos % 19 <= 15) || (add == 18 && pos % 19 >= 3) || (add == -20 && pos % 19 >= 3) || (add == 20 && pos % 19 <= 15) || add == 1 || add == -1 || add == 19 || add == -19 {
			if myMap[pos].Player != team {
				return val
			}
		}
	} else {
		return val
	}
	return checkLigne(myMap, pos + add, team, val + 1, add)
}

func CheckEnd(myMap []protocol.MapData, pos int, team int) bool {
	var nb int
	// horizontal
	nb = checkLigne(myMap, pos, team, 0, 1)
	nb += checkLigne(myMap, pos, team, 0, -1)
	if nb - 1 == 5 {
		//fmt.Printf("END 5 IN A ROW\n")
		return true
	}

	// vertical
	nb = checkLigne(myMap, pos, team, 0, 19)
	nb += checkLigne(myMap, pos, team, 0, -19)
	if nb - 1 == 5 {
		//fmt.Printf("END 5 IN A ROW\n")
		return true
	}

	// diagonal /
	nb = checkLigne(myMap, pos, team, 0, -18)
	nb += checkLigne(myMap, pos, team, 0, 18)
	if nb - 1 == 5 {
		//fmt.Printf("END 5 IN A ROW\n")
		return true
	}

	// diagonal \
	nb = checkLigne(myMap, pos, team, 0, -20)
	nb += checkLigne(myMap, pos, team, 0, 20)
	if nb - 1 == 5 {
		//fmt.Printf("END 5 IN A ROW\n")
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

func getNbPionTeamIndir(myMap []protocol.MapData, pos int, team int, dir int) int {
	//TODO: Real calcul for real diagonal and real line not a line in a border*
	if (pos + dir) < 0 || (pos + dir) >= (19 * 19) {
		return 0
	}
	if myMap[pos + dir].Player == team {
		return 1 + getNbPionTeamIndir(myMap, pos + dir, team, dir)
	}
	return 0
}

func Checkdoublethree(myMap []protocol.MapData, pos int, team int) bool {
	var nbline = 0
	var nbdiag = 0
	var pass = 0

	for i := 0; i < 8; i++ {
		if nbline < 2 {
			nbline = getNbPionTeamIndir(myMap, pos, team, Dirtab[i])
			if nbline >= 2 {
				pass = 1
			}
		}
		if pass == 0 && nbdiag < 2 {
			nbdiag = getNbPionTeamIndir(myMap, pos + (Dirtab[i] * 2), team, Dirtab[i])
		}
		if nbline == 2 && nbdiag == 2 {
			return true
		}
		pass = 0
	}
	return false
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
