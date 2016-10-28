package protocol

type MapData struct {
	Empty, Playable bool
	Team            int
}

const IDLE = "IDLE"
const START_OF_GAME = "START_OF_GAME"
const PLAY_TURN = "PLAY_TURN"
const END_OF_GAME = "END_OF_GAME"

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
