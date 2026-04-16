package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
	taskusecase "example.com/taskservice/internal/usecase/task"
)

// --- Task DTOs ---

type createTaskDTO struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      taskdomain.Status  `json:"status"`
	Recurrence  *recurrenceInputDTO `json:"recurrence,omitempty"`
}

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
}

type taskDTO struct {
	ID            int64              `json:"id"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Status        taskdomain.Status  `json:"status"`
	RecurrenceID  *int64             `json:"recurrence_id,omitempty"`
	ScheduledDate *time.Time         `json:"scheduled_date,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:            task.ID,
		Title:         task.Title,
		Description:   task.Description,
		Status:        task.Status,
		RecurrenceID:  task.RecurrenceID,
		ScheduledDate: task.ScheduledDate,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
	}
}

// --- Recurrence DTOs ---

type recurrenceInputDTO struct {
	Type      taskdomain.RecurrenceType `json:"type"`
	Interval  *int                      `json:"interval,omitempty"`
	MonthDays []int                     `json:"month_days,omitempty"`
	Dates     []time.Time               `json:"dates,omitempty"`
	Parity    *string                   `json:"parity,omitempty"`
}

func recurrenceInputToUsecase(dto recurrenceInputDTO) taskusecase.RecurrenceInput {
	return taskusecase.RecurrenceInput{
		Type:      dto.Type,
		Interval:  dto.Interval,
		MonthDays: dto.MonthDays,
		Dates:     dto.Dates,
		Parity:    dto.Parity,
	}
}