package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"websocket_manager/internal/model"
	"websocket_manager/internal/server/handlers"
	"websocket_manager/internal/session"
	"websocket_manager/internal/storage"
)

const (
	maxClients = 52
)

type Hub struct {
	context     context.Context
	connections map[uint64]*session.Session
	mu          *sync.Mutex
	storage     storage.Storage
	logger      *slog.Logger
}

func NewHub(context context.Context, storage storage.Storage, logger *slog.Logger) *Hub {
	return &Hub{
		context:     context,
		connections: make(map[uint64]*session.Session),
		mu:          &sync.Mutex{},
		storage:     storage,
		logger:      logger,
	}
}

func (h *Hub) Context() context.Context {
	return h.context
}

func (h *Hub) Logger() *slog.Logger {
	return h.logger
}

func (h *Hub) Register(session *session.Session) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.connections) >= maxClients {
		return errors.New("too many clients")
	}
	h.connections[session.ID()] = session
	return nil
}

func (h *Hub) Unregister(session *session.Session) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connections[session.ID()].Conn().Close()
	delete(h.connections, session.ID())
}

func (h *Hub) HandleMessage(msg *model.MessagePacketRequest) {
	h.Logger().Info("got message", "type", msg.MsgType, "from", msg.From, "to", msg.To, "msg", msg.Data)

	switch msg.MsgType {
	case model.SendMessage:
		ans := handlers.HandleSendMessage(h.storage, msg, h.logger.With("handler", "send_message", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
		// TODO: refactor probably
		if uow, err := h.storage.CreateUnitOfWork(); err != nil {
			users, err := uow.ChatRepository().GetAllUsersIDInChat(msg.To)
			if err != nil {
				return
			}
			for _, u := range users {
				if u == msg.From {
					continue
				}
				if _, ok := h.connections[u]; ok {
					getMessage := &model.MessagePacketRequest{MsgType: model.GetMessage, From: msg.From, To: msg.To, Data: json.RawMessage(fmt.Sprintf(`"%s"`, ans.Data))}
					h.connections[u].Enqueue(getMessage)
				}
			}
		}
	case model.UpdateMessage:
		ans := handlers.HandleUpdateMessage(h.storage, msg, h.logger.With("handler", "update_message", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	case model.DeleteMessage:
		ans := handlers.HandleDeleteMessage(h.storage, msg, h.logger.With("handler", "delete_message", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	case model.GetAllMessagesInChat: // TODO: should be limited to some reasonable amount
		ans := handlers.HandleGetAllMessagesInChat(h.storage, msg, h.logger.With("handler", "get_all_messages_in_chat", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	case model.CreateChat:
		ans := handlers.HandleCreateChat(h.storage, msg, h.logger.With("handler", "create_chat", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	case model.UpdateChat:
		ans := handlers.HandleUpdateChat(h.storage, msg, h.logger.With("handler", "update_chat", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	case model.DeleteChat:
		ans := handlers.HandleDeleteChat(h.storage, msg, h.logger.With("handler", "delete_chat", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	case model.AddUserToChat:
		ans := handlers.HandleAddUserToChat(h.storage, msg, h.logger.With("handler", "add_user_to_chat", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	case model.DeleteUserFromChat:
		ans := handlers.HandleDeleteUserFromChat(h.storage, msg, h.logger.With("handler", "delete_user_from_chat", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	case model.GetAllUsersIDInChat:
		ans := handlers.HandleGetLlUsersIDInChat(h.storage, msg, h.logger.With("handler", "get_all_users_id_in_chat", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	case model.GetAllUserChats:
		ans := handlers.HandleGetAllUserChats(h.storage, msg, h.logger.With("handler", "get_all_user_chats", "from", msg.From))
		h.connections[msg.From].Enqueue(ans)
	default:
		ans := &model.MessagePacketRequest{MsgType: model.SendMessage, From: 0, To: msg.From, Data: json.RawMessage("Internal Error")}
		h.connections[msg.From].Enqueue(ans)
	}
}
