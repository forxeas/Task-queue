CREATE TYPE jobs_status AS ENUM ('pending', 'failed', 'done', 'in_progressed');

CREATE TABLE IF NOT EXISTS jobs
(
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    status jobs_status DEFAULT 'pending',
    payload JSONB NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,
    available_at TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ
)