package lives

import (
	"fmt"
	"github.com/microservices-simulator-api/internal/boards"
	"github.com/microservices-simulator-api/internal/utils/json"
	"log"
	"sync"
)

// LiveBoard is a struct that holds the board and the websocket connections
type LiveBoard struct {
	sync.Mutex
	board         boards.Board
	state         *boards.BoardState
	activeClients map[int64]*ActiveClient
}

func NewLiveBoard(board boards.Board) *LiveBoard {
	return &LiveBoard{
		board: board,
		//state:         boards.BoardState{},
		activeClients: make(map[int64]*ActiveClient),
	}
}

func (b *LiveBoard) Broadcast(from int64, msg []byte) {
	for _, client := range b.activeClients {
		if client.Id != from {
			client.Write(msg)
		}
	}
}

func (b *LiveBoard) NewClientConnection(userId int64, channel Channel) error {
	fmt.Printf("New WebSocket connection established...\n")

	b.Lock()
	b.activeClients[userId] = NewActiveClient(userId, "Client 1", channel)
	b.Unlock()

	channel.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("Closing connection...\n")
		if client, ok := b.activeClients[userId]; ok {
			client.channel = nil
		}
		return nil
	})

	// Infinite loop
	for {
		// Read
		msg, err := channel.Receive()
		if err != nil {
			b.activeClients[userId] = &ActiveClient{}
			log.Println(err)
			return err
		}

		var inputEvent InputEvent
		if json.FromJson(msg, &inputEvent) != nil {
			log.Println(err)
		} else {
			outputEvent := OutputEvent{
				UserId:  userId,
				Event:   inputEvent.Event,
				Payload: inputEvent.Payload,
			}

			response, err := json.ToJson(outputEvent)
			if err != nil {
				log.Println(err)
			} else {
				b.Broadcast(userId, response)
			}
		}
	}
}
