package boards

type (
	BoardState struct {
		BoardId string   `json:"board_id"`
		Nodes   []Node   `json:"nodes"`
		Edges   []Edge   `json:"edges"`
		Packets []Packet `json:"packets"`
	}

	Node struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Position int    `json:"position"`
	}

	Edge struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Position int    `json:"position"`
	}

	Packet struct {
		Id       string   `json:"id"`
		Position Position `json:"position"`
		Content  Content  `json:"content"`
	}

	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	Content struct {
		Message string `json:"message"`
	}
)

func NewBoardState(boardId string) BoardState {
	return BoardState{
		BoardId: boardId,
	}
}
