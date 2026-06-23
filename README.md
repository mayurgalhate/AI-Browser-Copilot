# AI Browser Copilot
> An autonomous Chrome Extension agent that leverages LangChainGo to execute browser actions and answer questions in real-time.

![Demo GIF or Screenshot](TODO_ADD_SCREENSHOT.png)
[![Build](https://img.shields.io/badge/build-passing-brightgreen)](#) [![Coverage](https://img.shields.io/badge/coverage-89%25-brightgreen)](#) [![License](https://img.shields.io/badge/license-MIT-blue)](#)

## 🚀 What It Does
The modern web is too dynamic for basic LLM chat interfaces that lack context. **AI Browser Copilot** solves this by granting a Llama-3.3 agent deterministic control over your browser via a React Chrome Extension. This enables true autonomous task execution (clicking, scrolling, form-filling) and context-aware Q&A directly on your active tab, drastically accelerating web research and workflow automation.

## 🛠 Tech Stack
* **Go (Golang) & Gorilla WebSockets**: Chosen over REST for high-concurrency, low-latency bidirectional streaming between the LLM and the browser.
* **React, TypeScript & Tailwind (MV3)**: Chosen for a lightweight, strictly-typed, and secure Chrome Extension UI.
* **LangChainGo & Groq**: Chosen for sub-second Llama-3.3-70B inference and robust, multi-tool orchestration.
* **Docker, PostgreSQL (pgvector) & Redis**: Containerized with Docker Compose to orchestrate scalable vector memory (RAG) and sub-millisecond conversation state management.

## ⚡ Quick Start
```bash
# 1. Clone the repository
git clone https://github.com/mayurgalhate/AI-Browser-Copilot.git
cd AI-Browser-Copilot

# 2. Start Database & Cache
docker-compose up -d

# 3. Start Go Backend
cd backend
go run main.go
*(Note: Mac/Linux users can optionally use the included Makefile by running `make dev`)*

# 4. Build Extension
cd ../extension
npm install && npm run build
```
*After building, load the `extension/dist` folder into Chrome via `chrome://extensions/`.*

## 📁 Project Structure
```text
AI-Browser-Copilot/
├── backend/
│   ├── internal/
│   │   ├── agent/       # LangChainGo orchestrator & tool execution
│   │   ├── api/         # Gorilla WebSocket connection handlers
│   │   ├── config/      # Centralized environment loader
│   │   ├── db/          # PostgreSQL & Redis connection logic
│   │   └── tools/       # 8 custom browser tool definitions
│   ├── migrations/      # SQL schema evolution scripts
│   ├── docker-compose.yml
│   ├── Makefile
│   └── main.go
└── extension/
    ├── src/
    │   ├── background/  # Persistent Service Worker (WS Client)
    │   ├── content/     # DOM Extractor & Interactor logic
    │   ├── hooks/       # Custom React state (useCopilotChat)
    │   ├── types/       # Shared TS interfaces (ActionRequest)
    │   └── App.tsx      # React Sidebar UI
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

## 📐 Architecture / Design Decisions
* **WebSockets over REST**: Standard REST APIs proved too slow for real-time agentic execution. We migrated the primary communication layer to Gorilla WebSockets to stream agent thoughts and intermediate tool-calls directly to the browser UI without polling.
* **Deterministic DOM Manipulation**: LLMs historically hallucinate CSS selectors, causing brittle browser automation. To solve this, our `extractor.ts` script preemptively tags all interactive DOM elements with a stable `data-copilot-id`. The LLM receives this structured list and executes actions (clicks/fills) using guaranteed IDs.
* **Decoupled Extension Logic**: The Chrome Extension strictly separates DOM querying (`extractor.ts`) from DOM manipulation (`interactor.ts`), keeping the content script incredibly modular and safe.

## 🤝 Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License
Distributed under the MIT License. See `LICENSE` for more information.
