package task

import (
	"context"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
	CreateRecurrence(ctx context.Context, recurrence *taskdomain.Recurrence) (*taskdomain.Recurrence, error)
	GetActiveRecurrences(ctx context.Context) ([]taskdomain.Recurrence, error)
	DeactivateRecurrence(ctx context.Context, id int64) error
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
	CreateWithRecurrence(ctx context.Context, input CreateWithRecurrenceInput) (*taskdomain.Task, error)
	DeactivateRecurrence(ctx context.Context, id int64) error
	GenerateDailyTasks(ctx context.Context) error
}

type CreateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
}

type UpdateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
}

type CreateWithRecurrenceInput struct {
    Title       string
    Description string
    Status      taskdomain.Status
    Recurrence  RecurrenceInput
}

type RecurrenceInput struct {
    Type      taskdomain.RecurrenceType
    Interval  *int
    MonthDays []int
    Dates     []time.Time
    Parity    *string
}
