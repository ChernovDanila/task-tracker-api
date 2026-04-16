package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	RecurrenceID  *int64     `json:"recurrence_id,omitempty"`
    ScheduledDate *time.Time `json:"scheduled_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RecurrenceType string

const (
    RecurrenceDaily         RecurrenceType = "daily"
    RecurrenceMonthly       RecurrenceType = "monthly"
    RecurrenceSpecificDates RecurrenceType = "specific_dates"
    RecurrenceParity        RecurrenceType = "parity"
)

type Recurrence struct {
    ID          int64          `json:"id"`
    Title       string         `json:"title"`
    Description string         `json:"description"`
    Type        RecurrenceType `json:"type"`
    Interval    *int           `json:"interval,omitempty"`
    MonthDays   []int          `json:"month_days,omitempty"`
    Dates       []time.Time    `json:"dates,omitempty"`
    Parity      *string        `json:"parity,omitempty"`
    IsActive    bool           `json:"is_active"`
    CreatedAt   time.Time      `json:"created_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

func (r RecurrenceType) Valid() bool {
    switch r {
    case RecurrenceDaily, RecurrenceMonthly, RecurrenceSpecificDates, RecurrenceParity:
        return true
    default:
        return false
    }
}
