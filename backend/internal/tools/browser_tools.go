package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/tmc/langchaingo/tools"
)

// ActionRequest is sent to the browser
type ActionRequest struct {
	RequestID string `json:"request_id"`
	Action    string `json:"action"`
	Target    string `json:"target,omitempty"`
	Value     string `json:"value,omitempty"`
}

// ActionResponse is received from the browser
type ActionResponse struct {
	RequestID string `json:"request_id"`
	Result    string `json:"result"`
	Error     string `json:"error,omitempty"`
}

// SessionManager manages active WebSocket connections and pending tool actions
type SessionManager struct {
	mu             sync.Mutex
	PendingActions map[string]chan string
	SendToBrowser  func(sessionID string, msg []byte) error
}

var GlobalSessionManager = &SessionManager{
	PendingActions: make(map[string]chan string),
}

func (s *SessionManager) DispatchAction(sessionID, action, target, value string) (string, error) {
	requestID := fmt.Sprintf("%d", time.Now().UnixNano())
	ch := make(chan string, 1)

	s.mu.Lock()
	s.PendingActions[requestID] = ch
	s.mu.Unlock()

	req := ActionRequest{
		RequestID: requestID,
		Action:    action,
		Target:    target,
		Value:     value,
	}

	msg, _ := json.Marshal(req)

	if s.SendToBrowser != nil {
		if err := s.SendToBrowser(sessionID, msg); err != nil {
			s.mu.Lock()
			delete(s.PendingActions, requestID)
			s.mu.Unlock()
			return "", err
		}
	} else {
		return "", fmt.Errorf("no browser connection available")
	}

	// Wait for response or timeout
	select {
	case res := <-ch:
		return res, nil
	case <-time.After(30 * time.Second):
		s.mu.Lock()
		delete(s.PendingActions, requestID)
		s.mu.Unlock()
		return "", fmt.Errorf("timeout waiting for browser response")
	}
}

func (s *SessionManager) ResolveAction(requestID, result string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, ok := s.PendingActions[requestID]; ok {
		ch <- result
		delete(s.PendingActions, requestID)
	}
}

// BrowserTool is a generic tool that sends an action to the browser
type BrowserTool struct {
	NameDesc        string
	DescriptionDesc string
	SessionID       string
	ActionType      string
}

func (t *BrowserTool) Name() string        { return t.NameDesc }
func (t *BrowserTool) Description() string { return t.DescriptionDesc }
func (t *BrowserTool) Call(ctx context.Context, input string) (string, error) {
	// Simple parsing for tools that need target and value
	// For simplicity, input is treated as Target. Complex tools might need JSON input.
	return GlobalSessionManager.DispatchAction(t.SessionID, t.ActionType, input, "")
}

// CreateBrowserTools returns a list of tools for the LangChain agent
func CreateBrowserTools(sessionID string) []tools.Tool {
	return []tools.Tool{
		&BrowserTool{
			NameDesc:        "navigate",
			DescriptionDesc: "Useful for navigating to a specific URL. Input should be the full URL.",
			SessionID:       sessionID,
			ActionType:      "navigate",
		},
		&BrowserTool{
			NameDesc:        "extract_content",
			DescriptionDesc: "Reads all the text content of the current webpage. Use this when the user asks to summarize or answer questions about the current page. IMPORTANT: Input MUST be the exact word 'none'.",
			SessionID:       sessionID,
			ActionType:      "extract_content",
		},
		&BrowserTool{
			NameDesc:        "extract_elements",
			DescriptionDesc: "Extracts all interactive elements from the current page and returns a list with stable IDs. Call this before clicking or filling inputs. IMPORTANT: Input MUST be the exact word 'none'.",
			SessionID:       sessionID,
			ActionType:      "extract_elements",
		},
		&BrowserTool{
			NameDesc:        "click_element",
			DescriptionDesc: "Clicks an element on the page. Input must be the stable ID obtained from extract_elements.",
			SessionID:       sessionID,
			ActionType:      "click",
		},
		&BrowserTool{
			NameDesc:        "fill_input",
			DescriptionDesc: "Fills a text input. Input format should be a JSON string with 'id' (from extract_elements) and 'value'. Example: {\"id\":\"123\", \"value\":\"test\"}",
			SessionID:       sessionID,
			ActionType:      "fill_input",
		},
		&BrowserTool{
			NameDesc:        "scroll",
			DescriptionDesc: "Scrolls the page. Input MUST be one of: 'down', 'up', 'top', 'bottom'.",
			SessionID:       sessionID,
			ActionType:      "scroll",
		},
		&BrowserTool{
			NameDesc:        "browser_history",
			DescriptionDesc: "Navigates browser history or refreshes the page. Input MUST be one of: 'back', 'forward', 'refresh'.",
			SessionID:       sessionID,
			ActionType:      "browser_history",
		},
		&BrowserTool{
			NameDesc:        "get_current_url",
			DescriptionDesc: "Returns the current URL of the active tab. IMPORTANT: Input MUST be the exact word 'none'.",
			SessionID:       sessionID,
			ActionType:      "get_current_url",
		},
	}
}
