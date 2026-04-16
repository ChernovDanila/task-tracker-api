package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		UpdatedAt:   s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	return input, nil
}

func (s *Service) CreateWithRecurrence(ctx context.Context, input CreateWithRecurrenceInput) (*taskdomain.Task, error) {
    
    if input.Title == "" {
        return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
    }
    if !input.Recurrence.Type.Valid() {
        return nil, fmt.Errorf("%w: invalid recurrence type", ErrInvalidInput)
    }

    recurrence := &taskdomain.Recurrence{
        Type:      input.Recurrence.Type,
        Interval:  input.Recurrence.Interval,
        MonthDays: input.Recurrence.MonthDays,
        Dates:     input.Recurrence.Dates,
        Parity:    input.Recurrence.Parity,
        IsActive:  true,
        CreatedAt: s.now(),
    }
    created, err := s.repo.CreateRecurrence(ctx, recurrence)
    if err != nil {
        return nil, err
    }

    now := s.now()
    scheduledDate := now
    task := &taskdomain.Task{
        Title:         input.Title,
        Description:   input.Description,
        Status:        taskdomain.StatusNew,
        RecurrenceID:  &created.ID,
        ScheduledDate: &scheduledDate,
        CreatedAt:     now,
        UpdatedAt:     now,
    }
    return s.repo.Create(ctx, task)
}

func (s *Service) DeactivateRecurrence(ctx context.Context, id int64) error {
    if id <= 0 {
        return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
    }
    return s.repo.DeactivateRecurrence(ctx, id)
}

func (s *Service) GenerateDailyTasks(ctx context.Context) error {
    recurrences, err := s.repo.GetActiveRecurrences(ctx)
    if err != nil {
        return err
    }

    today := s.now()

    for _, rec := range recurrences {
        if !s.shouldCreateTask(rec, today) {
            continue
        }
        task := &taskdomain.Task{
    		Title:         rec.Title,
    		Description:   rec.Description,
    		Status:        taskdomain.StatusNew,
    		RecurrenceID:  &rec.ID,
    		ScheduledDate: &today,
    		CreatedAt:     today,
    		UpdatedAt:     today,
		}		
        _, err := s.repo.Create(ctx, task)
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *Service) shouldCreateTask(rec taskdomain.Recurrence, t time.Time) bool {
    day := t.Day()
    switch rec.Type {
    case taskdomain.RecurrenceDaily:
        if rec.Interval != nil && *rec.Interval > 1 {
            daysSince := int(t.Sub(rec.CreatedAt).Hours() / 24)
            return daysSince%*rec.Interval == 0
        }
        return true
    case taskdomain.RecurrenceMonthly:
        for _, d := range rec.MonthDays {
            if d == day {
                return true
            }
        }
    case taskdomain.RecurrenceSpecificDates:
        for _, d := range rec.Dates {
            if d.Year() == t.Year() && d.Month() == t.Month() && d.Day() == day {
                return true
            }
        }
    case taskdomain.RecurrenceParity:
        if rec.Parity != nil {
            if *rec.Parity == "even" && day%2 == 0 {
                return true
            }
            if *rec.Parity == "odd" && day%2 != 0 {
                return true
            }
        }
    }
    return false
}