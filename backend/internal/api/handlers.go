package api

import (
	"context"
	"net/http"

	"browser-copilot-backend/internal/agent"
	"browser-copilot-backend/internal/db"
	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
	Mode      string `json:"mode"` // "ask" or "act"
}

// HealthHandler returns API health status
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ChatHandler handles AI requests
func ChatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SessionID == "" {
		req.SessionID = "default_session"
	}

	// In Ask mode, we might do RAG first, but for now we just use the agent for both
	// The agent will decide whether to use tools based on its prompt

	res, err := agent.RunAgent(context.Background(), req.SessionID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": res})
}

// HistoryClearHandler clears the Redis session state
func HistoryClearHandler(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	err := db.RDB.Del(context.Background(), sessionID).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear history"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "history cleared", "session_id": sessionID})
}
