CREATE TABLE IF NOT EXISTS schedules (
    id            SERIAL PRIMARY KEY,
    client_name   VARCHAR(255) NOT NULL,
    start_at      TIMESTAMPTZ NOT NULL,
    end_at        TIMESTAMPTZ NOT NULL,
    location      VARCHAR(255) NOT NULL,
    latitude      DOUBLE PRECISION,
    longitude     DOUBLE PRECISION,
    clock_in_at   TIMESTAMPTZ,
    clock_in_lat  DOUBLE PRECISION,
    clock_in_lng  DOUBLE PRECISION,
    clock_out_at  TIMESTAMPTZ,
    clock_out_lat DOUBLE PRECISION,
    clock_out_lng DOUBLE PRECISION,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tasks (
    id            SERIAL PRIMARY KEY,
    schedule_id   INTEGER NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    title         VARCHAR(255) NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'pending',
    reason        TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_schedules_start_at ON schedules(start_at);
CREATE INDEX idx_tasks_schedule_id ON tasks(schedule_id);
