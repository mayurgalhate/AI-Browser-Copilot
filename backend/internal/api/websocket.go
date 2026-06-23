package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"browser-copilot-backend/internal/agent"
	brtools "browser-copilot-backend/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSClient represents a connected extension
type WSClient struct {
	SessionID string
	Conn      *websocket.Conn
	mu        sync.Mutex
}

var (
	clients   = make(map[string]*WSClient)
	clientsMu sync.Mutex
)

func init() {
	// Set the SendToBrowser function for the tools package
	brtools.GlobalSessionManager.SendToBrowser = func(sessionID string, msg []byte) error {
		clientsMu.Lock()
		client, ok := clients[sessionID]
		clientsMu.Unlock()

		if !ok {
			return nil
		}

		client.mu.Lock()
		defer client.mu.Unlock()
		return client.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

// WsHandler handles WebSocket connections
func WsHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrade WS", "error", err)
		return
	}

	// For simplicity, we use a single session ID or get from query param
	sessionID := c.Query("session_id")
	if sessionID == "" {
		sessionID = "default_session"
	}

	client := &WSClient{
		SessionID: sessionID,
		Conn:      ws,
	}

	clientsMu.Lock()
	clients[sessionID] = client
	clientsMu.Unlock()

	slog.Info("🟢 React Extension Connected!", "session_id", sessionID)

	defer func() {
		clientsMu.Lock()
		delete(clients, sessionID)
		clientsMu.Unlock()
		ws.Close()
		slog.Info("🔴 Extension Disconnected", "session_id", sessionID)
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		// Parse message from browser
		var resp brtools.ActionResponse
		if err := json.Unmarshal(message, &resp); err == nil && resp.RequestID != "" {
			// This is a tool response
			brtools.GlobalSessionManager.ResolveAction(resp.RequestID, resp.Result)
		} else {
			// Plain message, maybe chat or connection initialization
			msgStr := string(message)
			slog.Info("Received WS chat message", "msg", msgStr)
			
			go func(query string) {
				res, err := agent.RunAgent(context.Background(), sessionID, query)
				
				client.mu.Lock()
				defer client.mu.Unlock()
				
				if err != nil {
					slog.Error("Agent Error", "error", err)
					client.Conn.WriteMessage(websocket.TextMessage, []byte("⚠️ Error: "+err.Error()))
					return
				}
				client.Conn.WriteMessage(websocket.TextMessage, []byte(res))
			}(msgStr)
		}
	}
}
