package models

import "time"

// Conversation represents a user's session
type Conversation struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// Message represents a chat message
type Message struct {
	ID             int       `json:"id"`
	ConversationID string    `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// WebpageDocument represents a visited webpage context
type WebpageDocument struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// BrowserAction represents a tool invocation
type BrowserAction struct {
	ID         int       `json:"id"`
	SessionID  string    `json:"session_id"`
	ActionType string    `json:"action_type"`
	Payload    string    `json:"payload"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToolPayload is the structured message sent to the extension over WS
type ToolPayload struct {
	Tool   string `json:"tool"`
	Target string `json:"target,omitempty"`
	Value  string `json:"value,omitempty"`
}
