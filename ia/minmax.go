package ia

import (
	"fmt"
	"github.com/quentin-sommer/gomoku-back/protocol"
	"github.com/quentin-sommer/gomoku-back/referee"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

const (
	TWO_ALIGN   = 10
	THREE_ALIGN = 50
	FOUR_ALIGN  = 500
	// Compute : base + pawn taken
	BASE_PAWN_TAKEN = 300
	// Most important, wins over the rest every time
	FIVE_ALIGN = 10000
	MAX_INIT   = -420000
	MIN_INIT   = 420000
	INF_NEG    = MAX_INIT
	INF_POS    = MIN_INIT
)

var smallestIndex int
var highestIndex int

type MinMaxStruct struct {
	M      []protocol.MapData
	Player int8
	Depth  int
	End    bool
	Idx    int
}

func getOtherPlayer(player int8) int8 {
	if player == 0 {
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
	if cell.Empty {
		m[idx].Empty = false
		m[idx].Playable = false
		m[idx].Player = player

		return true
	}
	return false
}

func Eval(data *MinMaxStruct) int {
	ret := 0
	ret += FIVE_ALIGN * CountSequences(data, 5)
	ret += FOUR_ALIGN * CountSequences(data, 4)
	ret += THREE_ALIGN * CountSequences(data, 3)
	ret += TWO_ALIGN * CountSequences(data, 2)
	data.Player = getOtherPlayer(data.Player)
	ret -= FIVE_ALIGN * CountSequences(data, 5)
	ret -= FOUR_ALIGN * CountSequences(data, 4)
	ret -= THREE_ALIGN * CountSequences(data, 3)
	ret -= TWO_ALIGN * CountSequences(data, 2)
	data.Player = getOtherPlayer(data.Player)
	return ret
}

func caseNextToMe(m []protocol.MapData, idx int) bool {

	if m[idx].Player == -1 {
		if (idx-19 >= 0 && m[idx-19].Player != -1) || (idx+19 < protocol.MAP_SIZE && m[idx+19].Player != -1) ||
			(idx+1 < protocol.MAP_SIZE && m[idx+1].Player != -1) || (idx-1 >= 0 && m[idx-1].Player != -1) ||
			(idx-20 >= 0 && m[idx-20].Player != -1) || (idx+20 < protocol.MAP_SIZE && m[idx+20].Player != -1) ||
			(idx-18 >= 0 && m[idx-18].Player != -1) || (idx+18 < protocol.MAP_SIZE && m[idx+18].Player != -1) {
			return true
		}
	}
	return false
}

func max(data *MinMaxStruct, alpha, beta int) int {
	if data.Depth == 0 || data.End {
		return Eval(data)
	}
	max := MAX_INIT

	mapcp := make([]protocol.MapData, len(data.M))
	copy(mapcp, data.M)
	for i := 0; i < protocol.MAP_SIZE; i++ {
		if caseNextToMe(mapcp, i) {
			if playIdx(mapcp, i, data.Player) {
				captured, end, valid := referee.Exec(mapcp, i)
				if valid {
					tmp := min(&MinMaxStruct{mapcp, data.Player, data.Depth - 1, end, data.Idx}, alpha, beta)
					if captured > 0 {
						tmp += BASE_PAWN_TAKEN * (captured / 2)
					}

					if tmp > max {
						max = tmp
					}
					if max > alpha {
						alpha = max
					}
					if beta <= alpha {
						return max
					}
				}
				if captured > 0 {
					copy(mapcp, data.M)
				} else {
					mapcp[i] = data.M[i]
				}
			}
		}
	}
	return max
}

func min(data *MinMaxStruct, alpha, beta int) int {
	if data.Depth == 0 || data.End {
		return Eval(data)
	}
	min := MIN_INIT

	mapcp := make([]protocol.MapData, len(data.M))
	copy(mapcp, data.M)
	for i := 0; i < protocol.MAP_SIZE; i++ {
		if caseNextToMe(mapcp, i) {
			if playIdx(mapcp, i, getOtherPlayer(data.Player)) {
				captured, end, valid := referee.Exec(mapcp, i)
				if valid {
					tmp := max(&MinMaxStruct{mapcp, data.Player, data.Depth - 1, end, data.Idx}, alpha, beta)
					if captured > 0 {
						tmp -= BASE_PAWN_TAKEN * (captured / 2)
					}

					if tmp < min {
						min = tmp
					}
					if min < beta {
						beta = min
					}
					if beta <= alpha {
						return min
					}
				}
				if captured > 0 {
					copy(mapcp, data.M)
				} else {
					mapcp[i] = data.M[i]
				}
			}
		}
	}
	return min
}
func initSmallMax(m []protocol.MapData) {
	smallestIndex = -1
	for i := 0; i < protocol.MAP_SIZE; i++ {
		if !m[i].Empty {
			if smallestIndex == -1 {
				smallestIndex = i
			}
			highestIndex = i
		}
	}
	if smallestIndex >= 20 {
		smallestIndex -= 20
	} else {
		smallestIndex = 0
	}
	if (highestIndex + 21) < protocol.MAP_SIZE {
		highestIndex += 21
	} else {
		highestIndex = protocol.MAP_SIZE
	}
	// fmt.Println("iteration window size", highestIndex - smallestIndex, "(" + strconv.Itoa(smallestIndex) + "->" + strconv.Itoa(highestIndex) + ")")
}

func MinMaxBenchWrapper(m []protocol.MapData, player int8, depth int) int {
	f, err := os.Create("gomoku.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	return MinMax(m, player, depth)
}

func MinMax(m []protocol.MapData, player int8, depth int) int {
	start := time.Now().UnixNano()
	initSmallMax(m)
	maxval := MAX_INIT
	maxIdx := 0
	minval := MIN_INIT
	minIdx := 0

	alpha := INF_NEG
	beta := INF_POS

	mapcp := make([]protocol.MapData, len(m))
	copy(mapcp, m)
	for i := smallestIndex; i < highestIndex; i++ {
		if caseNextToMe(mapcp, i) {
			if playIdx(mapcp, i, player) {
				captured, end, valid := referee.Exec(mapcp, i)
				if valid {
					//tmp := Eval(&MinMaxStruct{mapcp, player, depth - 1, end, i})
					tmp := min(&MinMaxStruct{mapcp, player, depth - 1, end, i}, alpha, beta)
					if captured > 0 {
						tmp += BASE_PAWN_TAKEN * (captured / 2)
					}

					//fmt.Println(tmp, i)
					if tmp < minval {
						minval = tmp
						minIdx = i
					}
					if tmp > maxval {
						maxval = tmp
						maxIdx = i
					}
					if maxval > alpha {
						alpha = maxval
					}
				}
				if captured > 0 {
					copy(mapcp, m)
				} else {
					mapcp[i] = m[i]
				}
			}
		}
	}
	// fmt.Println("Id", maxIdx, maxval, minIdx, minval)
	/*if (minval < 0){
	  if (minval * -1 > maxval) {
	    return minIdx
	  }
	}*/
	end := time.Now().UnixNano()
	fmt.Println("Player", player, " depth", depth)
	fmt.Println("Real time taken by AI :", (end-start)/1000000, "ms")
	if minval <= (-FIVE_ALIGN + 1000) {
		return minIdx
	}
	return maxIdx
}
