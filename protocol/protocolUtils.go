package protocol

type MapData struct {
  Empty, Playable, LegitDoubleThree bool
  Player                            int8
}

const IDLE = "IDLE"
const START_OF_GAME = "START_OF_GAME"
const PLAY_TURN = "PLAY_TURN"
const END_OF_GAME = "END_OF_GAME"
const ENTER_ROOM = "ENTER_ROOM"
const REFRESH = "REFRESH"
const SET_AI_LEVEL = "SET_AI_LEVEL"
const MAP_SIZE = 19 * 19

type MessageIdle struct {
  Type string
}

type MessageStartOfGame struct {
  Type         string
  PlayerNumber int
}

type MessagePlayTurn struct {
  Type          string
  Map           []MapData
  TurnsPlayed   [2]int
  CapturedPawns [2]int
  IndexPlayed   int
}

type MessageEndOfGame struct {
  Type          string
  Map           []MapData
  TurnsPlayed   [2]int
  CapturedPawns [2]int
  Winner        int
}

type MessageSetAiLevel struct {
  Level int
}

type MessageEnterRoom struct {
  Type string
  Room int
	AiMode bool
}

type MessageRefresh struct {
  Type          string
  Map           []MapData
  TurnsPlayed   [2]int
  CapturedPawns [2]int
}

func SendEndOfGame(m []MapData, turnsPlayed [2]int, capturedPawns [2]int, winner int) *MessageEndOfGame {
  return &MessageEndOfGame{
    END_OF_GAME,
    m,
    turnsPlayed,
    capturedPawns,
    winner}
}

func SendPlayTurn(m []MapData, turnsPlayed [2]int, capturedPawns [2]int, indexPlayed int) *MessagePlayTurn {
  return &MessagePlayTurn{
    PLAY_TURN,
    m,
    turnsPlayed,
    capturedPawns,
    indexPlayed}
}

func SendStartOfGame(number int) *MessageStartOfGame {
  return &MessageStartOfGame{
    START_OF_GAME,
    number}
}

func SendIdle() *MessageIdle {
  return &MessageIdle{
    IDLE}
}

func SendRefresh(m []MapData, turnsPlayed [2]int, capturedPawns [2]int) *MessageRefresh {
  return &MessageRefresh{
    REFRESH,
    m,
    turnsPlayed,
    capturedPawns}
}

func InitGameData() ([]MapData, [2]int, [2]int) {
  myMap := make([]MapData, 19 * 19)
  for x := 0; x < 19 * 19; x++ {
    myMap[x].Empty = true
    myMap[x].Playable = true
    myMap[x].Player = -1
    myMap[x].LegitDoubleThree = false
  }
  var turnsPlayed [2]int
  turnsPlayed[0] = 0
  turnsPlayed[1] = 60

  var capturedPawns [2]int
  capturedPawns[0] = 0
  capturedPawns[1] = 0

  return myMap, turnsPlayed, capturedPawns
}

func IsInMap(myMap []MapData, x int, y int) bool {
  return x < 19 && x >= 0 && y < 19 && y >= 0
}
