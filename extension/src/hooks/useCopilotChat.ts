import { useEffect, useState } from 'react';
import type { ChatMessage } from '../types/messages';

declare const chrome: any;

export function useCopilotChat() {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [wsStatus, setWsStatus] = useState<'connected' | 'disconnected' | 'connecting'>('connecting');

  useEffect(() => {
    try {
      chrome.runtime.sendMessage({ type: "GET_WS_STATUS" }, (response: any) => {
        if (chrome.runtime.lastError) {
          console.warn("Background script not ready yet", chrome.runtime.lastError);
          return;
        }
        if (response && response.status) {
          setWsStatus(response.status);
        }
      });
    } catch (e) {
      console.warn("Message sending failed", e);
    }

    const messageListener = (message: any) => {
      if (message.type === "WS_STATUS") {
        setWsStatus(message.status);
        if (message.status === 'connected') {
          setMessages(prev => [...prev, { id: Date.now().toString(), sender: 'system', text: '⚡ Connected to AI Brain' }]);
        }
      } else if (message.type === "AI_MESSAGE") {
        setMessages(prev => [...prev, { id: Date.now().toString(), sender: 'ai', text: message.text }]);
      }
    };

    chrome.runtime.onMessage.addListener(messageListener);
    return () => chrome.runtime.onMessage.removeListener(messageListener);
  }, []);

  const addMessage = (message: ChatMessage) => {
    setMessages(prev => [...prev, message]);
  };

  return {
    messages,
    wsStatus,
    addMessage
  };
}
