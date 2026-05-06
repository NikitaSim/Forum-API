CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    password VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS threads (
    id BIGSERIAL PRIMARY KEY,
    author_id UUID REFERENCES users(id),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS thread_tags (
    thread_id BIGINT REFERENCES threads(id) ON DELETE CASCADE,
    tag VARCHAR(32),
    PRIMARY KEY (thread_id, tag)
);

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    author_id UUID REFERENCES users(id),
    thread_id BIGINT REFERENCES threads(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS attachments (
    id UUID PRIMARY KEY,
    thread_id BIGINT REFERENCES threads(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size BIGINT NOT NULL,
    path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS idempotency_keys (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    key VARCHAR(128) NOT NULL,
    request_hash TEXT,
    response_body JSONB,
    status_code INT,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (user_id, key)
);

CREATE INDEX idx_threads_author ON threads(author_id);
CREATE INDEX idx_posts_thread ON posts(thread_id);
CREATE INDEX idx_thread_tags_tag ON thread_tags(tag);