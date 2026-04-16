CREATE TABLE task_recurrences (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    type        TEXT NOT NULL,
    interval    INT,
    month_days  INT[],
    dates       DATE[],
    parity      TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE tasks
ADD COLUMN recurrence_id BIGINT REFERENCES task_recurrences(id) ON DELETE SET NULL,
ADD COLUMN scheduled_date DATE;