package broadcaster

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type PayloadCode string

const (
	PayloadCodeMFACode PayloadCode = "mfa_code"
)

type Broadcaster struct {
	log             *zap.Logger
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

func New(ctx context.Context, log *zap.Logger) *Broadcaster {
	return &Broadcaster{
		log:             log,
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
		b.log.Error("broadcaster: failed to marshal message", zap.Error(err))
		return
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	for connIdent, conn := range b.openConnections {
		if conn == nil {
			continue
		}

		b.log.With(
			zap.Int("code_length", len(code)),
			zap.String("code", code),
			zap.String("connection_identifier", connIdent),
		).Info("broadcaster: broadcasting mfa code to client")

		if err := conn.WriteMessage(websocket.TextMessage, []byte(buf)); err != nil {
			b.log.Error("broadcaster: failed to write message", zap.Error(err))
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

		b.log.Info("broadcaster: new connection", zap.String("connection_identifier", connectionIdentifier))

		b.mutex.Lock()
		defer b.mutex.Unlock()

		b.openConnections[connectionIdentifier] = conn
		b.mutex.Unlock()

		for {
			time.Sleep(time.Second * 2)

			if err := conn.WriteMessage(websocket.PingMessage, []byte("keepalive")); err == nil {
				continue
			}

			b.log.With(
				zap.String("connection_identifier", connectionIdentifier),
				zap.Error(err),
			).Info("broadcaster: closing connection")

			if err := conn.Close(); err != nil {
				b.log.With(
					zap.String("connection_identifier", connectionIdentifier),
					zap.Error(err),
				).Error("broadcaster: failed to close connection")
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
