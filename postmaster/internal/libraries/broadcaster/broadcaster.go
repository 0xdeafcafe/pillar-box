package broadcaster

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type PayloadCode string

const (
	PayloadCodeMFACode PayloadCode = "mfa_code"
)

type Broadcaster struct {
	mutex           sync.Mutex
	openConnections map[string]*websocket.Conn
	running         bool
}

type WebsocketMessage struct {
	Code    string                   `json:"code"`
	Payload *WebsocketMessagePayload `json:"payload"`
}

type WebsocketMessagePayload struct {
	MFACode *WebsocketMessagePayloadMFACode `json:"mfa_code"`
}

type WebsocketMessagePayloadMFACode struct {
	Code string `json:"code"`
}

// New creates a new Broadcaster instance. The Broadcaster is responsible for managing
// websocket connections and broadcasting messages to connected clients.
func New() *Broadcaster {
	return &Broadcaster{
		mutex:           sync.Mutex{},
		openConnections: make(map[string]*websocket.Conn),
		running:         false,
	}
}

func (b *Broadcaster) BroadcastMFACode(code string) {
	message := &WebsocketMessage{
		Code: string(PayloadCodeMFACode),
		Payload: &WebsocketMessagePayload{
			MFACode: &WebsocketMessagePayloadMFACode{
				Code: code,
			},
		},
	}

	buf, err := json.Marshal(message)
	if err != nil {
		log.Printf("broadcaster: failed to marshal message: %v", err)
		return
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	for connIdent, conn := range b.openConnections {
		if conn == nil {
			continue
		}

		log.Printf("broadcaster: failed to marshal message: %v code_length:%d code:%s connection_identifier:%s", err, len(code), code, connIdent)

		if err := conn.WriteMessage(websocket.TextMessage, []byte(buf)); err != nil {
			log.Printf("broadcaster: failed to write message: %v", err)
		}
	}
}

func (b *Broadcaster) ListenAndBroadcast() {
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(fmt.Sprintf("broadcaster: failed to upgrade connection: %v", err))
		}

		connectionIdentifier := uuid.New().String()

		log.Printf("broadcaster: new connection connection_identifier:%s", connectionIdentifier)

		b.mutex.Lock()
		defer b.mutex.Unlock()

		b.openConnections[connectionIdentifier] = conn
		b.mutex.Unlock()

		for {
			time.Sleep(time.Second * 2)

			if err := conn.WriteMessage(websocket.PingMessage, []byte("keepalive")); err == nil {
				continue
			}

			log.Printf("broadcaster: closing connection: %v connection_identifier:%s", err, connectionIdentifier)

			if err := conn.Close(); err != nil {
				log.Printf("broadcaster: failed to close connection: %v connection_identifier:%s", err, connectionIdentifier)
			}

			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.openConnections[connectionIdentifier] = nil
		}
	})

	// TODO(afr): make this configurable
	go http.ListenAndServe(":3500", nil)

	b.running = true
}
