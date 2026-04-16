package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, description, status, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
	created, err := scanTask(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING id, title, description, status, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.UpdatedAt, task.ID)
	updated, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		ORDER BY id DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task   taskdomain.Task
		status string
	)

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)

	return &task, nil
}

func (r *Repository) CreateRecurrence(ctx context.Context, rec *taskdomain.Recurrence) (*taskdomain.Recurrence, error) {
    const query = `
        INSERT INTO task_recurrences (title, description, type, interval, month_days, dates, parity, is_active, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, title, description, type, interval, month_days, dates, parity, is_active, created_at
    `
    row := r.pool.QueryRow(ctx, query,
        rec.Title,
        rec.Description,
        rec.Type,
        rec.Interval,
        rec.MonthDays,
        rec.Dates,
        rec.Parity,
        rec.IsActive,
        rec.CreatedAt,
    )
    return scanRecurrence(row)
}

func (r *Repository) GetActiveRecurrences(ctx context.Context) ([]taskdomain.Recurrence, error) {
    const query = `
        SELECT id, title, description, type, interval, month_days, dates, parity, is_active, created_at
        FROM task_recurrences
        WHERE is_active = true
    `
    rows, err := r.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []taskdomain.Recurrence
    for rows.Next() {
        rec, err := scanRecurrence(rows)
        if err != nil {
            return nil, err
        }
        result = append(result, *rec)
    }
    return result, rows.Err()
}

func (r *Repository) DeactivateRecurrence(ctx context.Context, id int64) error {
    const query = `UPDATE task_recurrences SET is_active = false WHERE id = $1`
    result, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return err
    }
    if result.RowsAffected() == 0 {
        return taskdomain.ErrNotFound
    }
    return nil
}

func scanRecurrence(scanner taskScanner) (*taskdomain.Recurrence, error) {
    var rec taskdomain.Recurrence
    var recType string
    err := scanner.Scan(
        &rec.ID,
        &rec.Title,
        &rec.Description,
        &recType,
        &rec.Interval,
        &rec.MonthDays,
        &rec.Dates,
        &rec.Parity,
        &rec.IsActive,
        &rec.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    rec.Type = taskdomain.RecurrenceType(recType)
    return &rec, nil
}