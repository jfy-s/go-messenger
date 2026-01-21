package session

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"websocket_manager/internal/jwt"
	"websocket_manager/internal/model"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

type Hub interface {
	Context() context.Context
	Logger() *slog.Logger
	Register(session *Session) error
	Unregister(session *Session)
	HandleMessage(msg *model.MessagePacketRequest)
}

type Session struct {
	hub  Hub
	conn *websocket.Conn
	id   uint64
	send chan []byte
}

func (s *Session) ID() uint64 {
	return s.id
}

func (s *Session) Conn() *websocket.Conn {
	return s.conn
}

func (s *Session) Enqueue(msgPkt *model.MessagePacketRequest) {
	bytes, err := msgPkt.ToBytes()
	if err != nil {
		s.hub.Logger().Error("failed to convert message packet to bytes", "error", err)
		return
	}
	s.send <- bytes
}

func (s *Session) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		s.conn.Close()
	}()
	for {
		select {
		case <-s.hub.Context().Done():
			s.hub.Logger().Info("context closed")
			return

		case message, ok := <-s.send:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				s.hub.Logger().Error("failed to write message", "error", "channel closed")
				return
			}

			w, err := s.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				s.hub.Logger().Error("failed to get next writer", "error", err)
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				s.hub.Logger().Error("failed to close writer", "error", err)
				return
			}
		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}

func (s *Session) readPump() {
	defer func() {
		s.hub.Unregister(s)
		s.conn.Close()
		close(s.send)
	}()

	s.conn.SetReadLimit(maxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(func(string) error { s.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		select {
		case <-s.hub.Context().Done():
			s.hub.Logger().Info("context closed")
			return
		default:
			_, message, err := s.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					s.hub.Logger().Error("unexpected close", "error", err)
				}
				break
			}
			MessagePacketRequest, err := model.ByteToMessagePacketRequest(message)
			if err != nil {
				s.hub.Logger().Error("failed to convert message to message packet", "error", err)
				continue
			}
			MessagePacketRequest.From = s.id
			s.hub.HandleMessage(MessagePacketRequest)
		}
	}
}

func ServeWs(hub Hub, w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	tokenStr = tokenStr[len("Bearer "):]
	id, err := jwt.ParseToken(tokenStr)
	if err != nil {
		hub.Logger().Error("failed to parse token", "error", err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		hub.Logger().Error("failed to upgrade connection", "error", err)
		return
	}
	session := &Session{hub: hub, conn: conn, id: id, send: make(chan []byte)}

	if err := session.hub.Register(session); err != nil {
		hub.Logger().Error("failed to register connection", "error", err)
		return
	}

	go session.writePump()
	go session.readPump()
}
