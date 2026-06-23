import { useState, useRef, useEffect } from 'react';
import { useCopilotChat } from './hooks/useCopilotChat';

declare const chrome: any;

function App() {
  const { messages, wsStatus, addMessage } = useCopilotChat();
  const [input, setInput] = useState('');
  const [mode, setMode] = useState<'ask' | 'act'>('ask');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const sendMessage = () => {
    if (input.trim() === '') return;
    
    chrome.tabs.query({ active: true, currentWindow: true }, (tabs: any[]) => {
      const currentTab = tabs[0];
      const contextStr = currentTab ? `\n\n[Context: User is currently on URL: ${currentTab.url} with Title: "${currentTab.title}"]` : '';
      
      let payload = input + contextStr;
      if (mode === 'act') {
        payload = `[ACT MODE] ${payload}`;
      } else {
        payload = `[ASK MODE] ${payload}`;
      }

      addMessage({ id: Date.now().toString(), sender: 'user', text: input });
      setInput('');

      try {
        chrome.runtime.sendMessage({ type: "SEND_CHAT", text: payload }, (response: any) => {
          if (chrome.runtime.lastError) {
            console.warn("Background script not ready yet", chrome.runtime.lastError);
            addMessage({ id: Date.now().toString(), sender: 'system', text: '⚠️ Connection Error: Background script not ready.' });
            return;
          }
          if (!response || !response.success) {
            addMessage({ id: Date.now().toString(), sender: 'system', text: '⚠️ Connection Error: ' + (response?.error || 'Unknown') });
          }
        });
      } catch (e) {
        console.warn("Message sending failed", e);
        addMessage({ id: Date.now().toString(), sender: 'system', text: '⚠️ Extension Context Invalidated. Please reload.' });
      }
    });
  };

  return (
    <div className="flex flex-col h-screen w-full font-sans select-none overflow-hidden text-slate-200">
      
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-white/10 bg-white/5 backdrop-blur-md">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-indigo-500 to-purple-500 flex items-center justify-center shadow-lg shadow-indigo-500/30">
            <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div>
            <h1 className="text-lg font-bold tracking-tight text-white">Copilot</h1>
            <p className="text-[10px] uppercase tracking-wider text-indigo-300 font-semibold opacity-80">
              {wsStatus === 'connected' ? '🟢 Online' : '🔴 Offline'}
            </p>
          </div>
        </div>

        {/* Mode Toggle */}
        <div className="flex bg-black/40 rounded-full p-1 border border-white/5">
          <button 
            onClick={() => setMode('ask')}
            className={`px-3 py-1 text-xs font-medium rounded-full transition-all ${mode === 'ask' ? 'bg-indigo-500 text-white shadow-md shadow-indigo-500/20' : 'text-slate-400 hover:text-white'}`}
          >
            Ask
          </button>
          <button 
            onClick={() => setMode('act')}
            className={`px-3 py-1 text-xs font-medium rounded-full transition-all ${mode === 'act' ? 'bg-purple-500 text-white shadow-md shadow-purple-500/20' : 'text-slate-400 hover:text-white'}`}
          >
            Act
          </button>
        </div>
      </div>
      
      {/* Messages Area */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.length === 0 && (
          <div className="h-full flex flex-col items-center justify-center text-center opacity-50 space-y-3">
            <div className="w-16 h-16 rounded-full bg-white/5 flex items-center justify-center mb-2">
              <span className="text-2xl">✨</span>
            </div>
            <p className="text-sm max-w-[200px]">How can I help you navigate or understand the web today?</p>
          </div>
        )}
        
        {messages.map((msg) => {
          const isUser = msg.sender === 'user';
          const isSystem = msg.sender === 'system';
          
          if (isSystem) {
            return (
              <div key={msg.id} className="flex justify-center my-2">
                <span className="text-[11px] font-mono tracking-wider text-slate-500 bg-black/30 px-3 py-1 rounded-full border border-white/5">
                  {msg.text}
                </span>
              </div>
            );
          }

          return (
            <div key={msg.id} className={`flex ${isUser ? 'justify-end' : 'justify-start'}`}>
              {!isUser && (
                <div className="w-6 h-6 rounded-full bg-gradient-to-tr from-indigo-500 to-purple-500 mt-1 mr-2 flex-shrink-0 flex items-center justify-center">
                  <span className="text-[10px]">AI</span>
                </div>
              )}
              <div className={`max-w-[85%] rounded-2xl px-4 py-2.5 text-[13px] leading-relaxed shadow-sm ${
                isUser 
                  ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-tr-sm border border-blue-500/50' 
                  : 'bg-white/10 text-slate-100 rounded-tl-sm border border-white/10 backdrop-blur-sm whitespace-pre-wrap'
              }`}>
                {msg.text}
              </div>
            </div>
          );
        })}
        <div ref={messagesEndRef} />
      </div>

      {/* Input Area */}
      <div className="p-4 border-t border-white/10 bg-black/20">
        <div className="flex items-center gap-2 bg-white/5 border border-white/10 rounded-2xl p-1.5 focus-within:border-indigo-500/50 focus-within:bg-white/10 transition-all shadow-inner shadow-black/20">
          <input 
            type="text" 
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
            placeholder={mode === 'ask' ? "Ask about this page..." : "Tell me what to do..."}
            className="flex-1 bg-transparent border-none outline-none text-[13px] px-3 py-1.5 text-white placeholder:text-slate-500"
            disabled={wsStatus !== 'connected'}
          />
          <button 
            onClick={sendMessage} 
            disabled={!input.trim() || wsStatus !== 'connected'}
            className="w-8 h-8 rounded-full bg-indigo-500 text-white flex items-center justify-center hover:bg-indigo-400 disabled:opacity-50 disabled:hover:bg-indigo-500 transition-colors flex-shrink-0 shadow-lg shadow-indigo-500/20"
          >
            <svg className="w-4 h-4 ml-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}

export default App;