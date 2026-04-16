ALTER TABLE tasks
DROP COLUMN IF EXISTS scheduled_date,
DROP COLUMN IF EXISTS recurrence_id;

DROP TABLE IF EXISTS task_recurrences;