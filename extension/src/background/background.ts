// Tells TypeScript about the global 'chrome' variable
declare const chrome: any;

let ws: WebSocket | null = null;
const sessionID = "session_" + Date.now();

// Setup side panel behavior
chrome.sidePanel.setPanelBehavior({ openPanelOnActionClick: true }).catch(console.error);

function connectWebSocket() {
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
    return;
  }

  ws = new WebSocket(`ws://localhost:8080/ws?session_id=${sessionID}`);

  ws.onopen = () => {
    console.log("🟢 WS Connected to Go Backend");
    chrome.runtime.sendMessage({ type: "WS_STATUS", status: "connected" }).catch(() => {});
  };

  ws.onmessage = async (event) => {
    try {
      // Check if it's a tool action request from Go
      const data = JSON.parse(event.data);
      if (data.request_id && data.action) {
        handleToolAction(data);
        return;
      }
    } catch (e) {
      // It's a plain chat message from the AI
      chrome.runtime.sendMessage({ type: "AI_MESSAGE", text: event.data }).catch(() => {});
    }
  };

  ws.onclose = () => {
    console.log("🔴 WS Disconnected, reconnecting...");
    chrome.runtime.sendMessage({ type: "WS_STATUS", status: "disconnected" }).catch(() => {});
    setTimeout(connectWebSocket, 3000);
  };

  ws.onerror = (err: any) => {
    console.error("WS Error:", err);
  };
}

async function handleToolAction(data: any) {
  try {
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
    if (!tab || !tab.id) {
      sendToolResponse(data.request_id, "Error: No active tab found.");
      return;
    }

    if (data.action === "navigate") {
      await chrome.tabs.update(tab.id, { url: data.target });
      sendToolResponse(data.request_id, "Navigation successful");
      return;
    }

    // Forward other actions to the content script
    const response = await chrome.tabs.sendMessage(tab.id, {
      type: "EXECUTE_ACTION",
      action: data.action,
      target: data.target,
      value: data.value
    });

    sendToolResponse(data.request_id, response.result || "Action executed");
  } catch (err: any) {
    sendToolResponse(data.request_id, `Error: ${err.message}`);
  }
}

function sendToolResponse(requestID: string, result: string) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ request_id: requestID, result }));
  }
}

// Listen for messages from the Side Panel
chrome.runtime.onMessage.addListener((message: any, _sender: any, sendResponse: any) => {
  if (message.type === "SEND_CHAT") {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(message.text);
      sendResponse({ success: true });
    } else {
      sendResponse({ success: false, error: "WS not connected" });
    }
  } else if (message.type === "GET_WS_STATUS") {
    sendResponse({ status: ws?.readyState === WebSocket.OPEN ? "connected" : "disconnected" });
  }
});

// Initialize connection
connectWebSocket();
