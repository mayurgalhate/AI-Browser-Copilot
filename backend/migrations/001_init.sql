-- Enable the pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Conversations table to group messages
CREATE TABLE IF NOT EXISTS conversations (
    id VARCHAR(50) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB DEFAULT '{}'::jsonb
);

-- Messages table for chat history
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    conversation_id VARCHAR(50) REFERENCES conversations(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Store visited webpages
CREATE TABLE IF NOT EXISTS webpage_documents (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    title TEXT,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(url)
);

-- Store embeddings for webpage chunks
CREATE TABLE IF NOT EXISTS embeddings (
    id SERIAL PRIMARY KEY,
    document_id INTEGER REFERENCES webpage_documents(id) ON DELETE CASCADE,
    chunk_text TEXT NOT NULL,
    embedding vector(768), -- Adjust dimensionality based on the embedding model (e.g. 768 for nomic-embed-text)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create HNSW index for fast similarity search
CREATE INDEX ON embeddings USING hnsw (embedding vector_l2_ops);

-- Log of browser actions taken by the agent
CREATE TABLE IF NOT EXISTS browser_actions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(50) NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    payload JSONB DEFAULT '{}'::jsonb,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
