# AI Browser Copilot
> An autonomous Chrome Extension agent that leverages LangChainGo to execute browser actions and answer questions in real-time.


## 🚀 What It Does
The modern web is too dynamic for basic LLM chat interfaces that lack context. **AI Browser Copilot** solves this by granting a Llama-3.3 agent deterministic control over your browser via a React Chrome Extension. This enables true autonomous task execution (clicking, scrolling, form-filling) and context-aware Q&A directly on your active tab, drastically accelerating web research and workflow automation.


## 🛠 Tech Stack

| Technology | Role | Why Chosen |
| :--- | :--- | :--- |
| **Go (Golang)** | Backend Server | High-concurrency, low-latency performance at scale |
| **Gorilla WebSockets** | Real-time Transport | Bidirectional streaming between LLM and browser — REST can't match it |
| **React + TypeScript** | Extension UI | Lightweight, strictly-typed, secure Chrome Extension (MV3) |
| **Tailwind CSS** | Styling | Utility-first, minimal bundle size for extension constraints |
| **LangChainGo** | Agent Orchestration | Robust multi-tool orchestration natively in Go |
| **Groq + Llama-3.3-70B** | LLM Inference | Sub-second inference speed |
| **PostgreSQL + pgvector** | Vector Memory | Scalable RAG / semantic search storage |
| **Redis** | Session State | Sub-millisecond conversation context retrieval |
| **Docker + Compose** | Containerization | Orchestrates all services in a single reproducible environment |


## ⚡ Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/mayurgalhate/AI-Browser-Copilot.git
cd AI-Browser-Copilot
```

```bash
# 2. Start Database & Cache
docker-compose up -d
```

```bash
# 3. Start Go Backend
cd backend
go run main.go
```
> 💡 **Mac/Linux:** Optionally run `make dev` instead of `go run main.go`

```bash
# 4. Build Chrome Extension
cd ../extension
npm install && npm run build
```

> 🔌 **Load Extension:** Open `chrome://extensions/` → Enable **Developer Mode** → Click **Load Unpacked** → Select `extension/dist/`


## 📁 Project Structure

```text
AI-Browser-Copilot/
│
├── backend/
│   ├── internal/
│   │   ├── agent/          # LangChainGo orchestrator & tool execution
│   │   ├── api/            # Gorilla WebSocket connection handlers
│   │   ├── config/         # Centralized environment loader
│   │   ├── db/             # PostgreSQL & Redis connection logic
│   │   └── tools/          # 8 custom browser tool definitions
│   ├── migrations/         # SQL schema evolution scripts
│   ├── docker-compose.yml
│   ├── Makefile
│   └── main.go
│
└── extension/
    ├── src/
    │   ├── background/     # Persistent Service Worker (WebSocket Client)
    │   ├── content/        # DOM Extractor & Interactor logic
    │   ├── hooks/          # Custom React state (useCopilotChat)
    │   ├── types/          # Shared TS interfaces (ActionRequest)
    │   └── App.tsx         # React Sidebar UI
    └── vite.config.ts
```

## 🔧 Configuration
Create a `.env` file in the `backend/` directory based on `.env.example`:

| Environment Variable | Description | Default / Example |
| :--- | :--- | :--- |
| `GROQ_API_KEY` | Required API key for Llama-3.3 inference | `gsk_your_key_here` |
| `GROQ_BASE_URL` | Base URL for OpenAI-compatible Groq endpoint | `https://api.groq.com/openai/v1` |
| `DB_URL` | PostgreSQL connection string | `postgres://user:password@localhost...` |
| `REDIS_URL` | Redis connection address | `localhost:6379` |

## 🔌 Core API Endpoints
The Go backend exposes the following endpoints:

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/health` | `GET` | Checks active connections to PostgreSQL and Redis. |
| `/ws` | `GET` | Opens a persistent WebSocket connection for real-time agent streaming. |
| `/chat` | `POST` | Stateless REST fallback for LLM inference. |
| `/history/:id` | `DELETE` | Clears Redis conversation state for a given session. `:id` = unique session UUID. |

## 🧪 Running Tests
```bash
cd backend
make test
```
**What's covered:** Currently covers agent initialization, environment validation, and LLM provider connectivity ensuring the core LangChainGo pipeline boots successfully.


## 🏗 Architecture

```text
┌─────────────────────────────────────────────────────┐
│                  Chrome Extension                    │
│                                                      │
│   ┌─────────────┐        ┌──────────────────────┐   │
│   │  React UI   │◄──────►│  Service Worker (BG) │   │
│   │ (Sidebar)   │        │   WebSocket Client   │   │
│   └─────────────┘        └──────────┬───────────┘   │
│                                     │               │
│   ┌─────────────────────────────┐   │               │
│   │     Content Script          │   │               │
│   │  DOM Extractor & Interactor │   │               │
│   └─────────────────────────────┘   │               │
└─────────────────────────────────────┼───────────────┘
                                      │ WebSocket /ws
                        ┌─────────────▼───────────────┐
                        │        Go Backend            │
                        │                             │
                        │  ┌────────────────────────┐ │
                        │  │   Gorilla WS Handler   │ │
                        │  └──────────┬─────────────┘ │
                        │             │               │
                        │  ┌──────────▼─────────────┐ │
                        │  │   LangChainGo Agent    │ │
                        │  │   + Tool Orchestrator  │ │
                        │  └──────────┬─────────────┘ │
                        │             │               │
                        │      ┌──────▼──────┐        │
                        │      │  Groq API   │        │
                        │      │ Llama-3.3   │        │
                        │      └─────────────┘        │
                        └─────────────────────────────┘
                                      │
                 ┌────────────────────┼─────────────────────┐
                 │                    │                     │
        ┌────────▼──────┐   ┌─────────▼──────┐   ┌────────▼──────┐
        │  PostgreSQL   │   │     Redis      │   │   pgvector    │
        │  (Persistent) │   │  (Session CTX) │   │  (RAG/Memory) │
        └───────────────┘   └────────────────┘   └───────────────┘
```



 ## 🔄 How It Works
 
1. User types a prompt in the React Sidebar
2. Service Worker sends it via WebSocket to Go backend `/ws`
3. LangChainGo Agent receives prompt + fetches Redis session context
4. Agent calls Groq (Llama-3.3-70B) for reasoning
5. If a browser action is needed → Agent triggers one of 8 custom tools
6. Tool response sent back to Content Script → executes DOM action
7. Result streamed back to Sidebar in real-time
8. Conversation state saved to Redis + vector stored in pgvector


## 🧰 Agent Tools

| # | Tool | What It Does |
| :--- | :--- | :--- |
| 1 | `get_page_content` | Extracts full DOM text of current tab |
| 2 | `click_element` | Clicks a DOM element by selector |
| 3 | `fill_input` | Fills an input field with given value |
| 4 | `scroll_page` | Scrolls up/down on the page |
| 5 | `get_current_url` | Returns the active tab URL |
| 6 | `navigate_to` | Redirects browser to a given URL |
| 7 | `summarize_page` | Summarizes current page via LLM |
| 8 | `search_memory` | Queries pgvector for past context (RAG) |

---

## 🌊 Data Flow

```text
User Prompt
    │
    ▼
Service Worker (WS Client)
    │
    ▼
Go Backend → Redis (fetch session)
    │
    ▼
LangChainGo Agent → Groq Llama-3.3-70B
    │
    ├──► Browser Tool → Content Script → DOM
    │
    ▼
Stream Response → React Sidebar
    │
    ▼
Save → Redis (state) + pgvector (memory)
```
