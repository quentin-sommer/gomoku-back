package protocol

type MapData struct {
	Empty, Playable bool
	Player          int
}

const IDLE = "IDLE"
const START_OF_GAME = "START_OF_GAME"
const PLAY_TURN = "PLAY_TURN"
const END_OF_GAME = "END_OF_GAME"
const ENTER_ROOM = "ENTER_ROOM"
const REFRESH = "REFRESH"

type MessageIdle struct {
	Type string
}

type MessageStartOfGame struct {
	Type         string
	PlayerNumber int
}

type MessagePlayTurn struct {
	Type string
	Map  []MapData
}

type MessageEndOfGame struct {
	Type   string
	Winner int
	Map    []MapData
}

type MessageEnterRoom struct {
	Type	string
	Room	int
}

type MessageRefresh struct {
	Type  string
	Map   []MapData
}

func SendEndOfGame(m []MapData, winner int) (*MessageEndOfGame) {
	return &MessageEndOfGame{
		END_OF_GAME,
		winner,
		m}
}

func SendPlayTurn(m []MapData) (*MessagePlayTurn) {
	return &MessagePlayTurn{
		PLAY_TURN,
		m}
}

func SendStartOfGame(number int) (*MessageStartOfGame) {
	return &MessageStartOfGame{
		START_OF_GAME,
		number}
}

func SendIdle() (*MessageIdle) {
	return &MessageIdle{
		IDLE}
}

func SendRefresh(m []MapData) (*MessageRefresh) {
	return &MessageRefresh{
		REFRESH,
		m}
}
func InitMap() ([]MapData) {
	myMap := make([]MapData, 19 * 19)
	for x := 0; x < 19 * 19; x++ {
		myMap[x].Empty = true
		myMap[x].Playable = true
		myMap[x].Player = -1
	}
	return myMap
}
