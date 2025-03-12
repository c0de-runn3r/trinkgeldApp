package models

import (
	"encoding/json"
	"time"
)

// Location represents the 'locations' collection
type Location struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Worker represents the 'workers' collection
type Worker struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WorkerLocation represents the 'worker_locations' collection
type WorkerLocation struct {
	ID         string    `json:"id"`
	WorkerID   string    `json:"worker_id"`
	LocationID string    `json:"location_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// WorkShift represents the 'work_shifts' collection
type WorkShift struct {
	ID          string    `json:"id"`
	WorkerID    string    `json:"worker_id"`
	LocationID  string    `json:"location_id"`
	Date        time.Time `json:"date"`
	HoursWorked float64   `json:"hours_worked"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DailyTip represents the 'daily_tips' collection
type DailyTip struct {
	ID         string    `json:"id"`
	LocationID string    `json:"location_id"`
	Date       time.Time `json:"date"`
	TotalTips  float64   `json:"total_tips"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// WorkerTip represents the 'worker_tips' collection
type WorkerTip struct {
	ID          string    `json:"id"`
	WorkerID    string    `json:"worker_id"`
	LocationID  string    `json:"location_id"`
	Date        time.Time `json:"date"`
	HoursWorked float64   `json:"hours_worked"`
	TipsEarned  float64   `json:"tips_earned"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MonthlySummary represents the 'monthly_summaries' collection
type MonthlySummary struct {
	ID         string    `json:"id"`
	WorkerID   string    `json:"worker_id"`
	LocationID string    `json:"location_id"`
	Month      string    `json:"month"`
	TotalHours float64   `json:"total_hours"`
	TotalTips  float64   `json:"total_tips"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UnmarshalJSON customizes the JSON unmarshaling of DailyTip
func (d *DailyTip) UnmarshalJSON(data []byte) error {
	type Alias DailyTip
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the date string to time.Time, handling the specific format
	parsedDate, err := time.Parse("2006-01-02 15:04:05.000Z", aux.Date)
	if err != nil {
		return err
	}
	d.Date = parsedDate
	return nil
}

// UnmarshalJSON customizes the JSON unmarshaling of WorkShift
func (w *WorkShift) UnmarshalJSON(data []byte) error {
	type Alias WorkShift
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(w),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the date string to time.Time, handling the specific format
	parsedDate, err := time.Parse("2006-01-02 15:04:05.000Z", aux.Date)
	if err != nil {
		return err
	}
	w.Date = parsedDate
	return nil
}
