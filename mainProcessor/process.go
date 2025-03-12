package mainprocessor

import (
	"log"
	"math"
	"time"
	"trinkgeldApp/models"
)

// This function should calculate the tips earned by each worker per day
// based on the hours worked and the total tips of the location
// It should return a slice of WorkerTip structs
// Each WorkerTip struct should contain the worker ID, location ID, date, hours worked, tips earned, and timestamps
// The tips earned should be calculated as follows:
// tips earned = (hours worked / total hours worked) * total tips
// where total hours worked is the sum of hours worked by all workers at the location on that day
// and total tips is the total tips of the location on that day
func CalculateTipAmountsPerWorkerPerDay(shifts []*models.WorkShift, dailyTips []*models.DailyTip) ([]*models.WorkerTip, error) {
	tipsEarned := make([]*models.WorkerTip, 0)

	// Group shifts by date
	shiftsByDate := make(map[time.Time][]*models.WorkShift)
	for _, shift := range shifts {
		date := shift.Date.Truncate(24 * time.Hour)
		shiftsByDate[date] = append(shiftsByDate[date], shift)
	}

	// Iterate through each day
	for date, dailyShifts := range shiftsByDate {
		var dailyTip *models.DailyTip
		for _, tip := range dailyTips {
			if tip.Date.Truncate(24 * time.Hour).Equal(date) {
				dailyTip = tip
				break
			}
		}
		if dailyTip == nil {
			continue
		}

		totalHours := 0.0
		for _, shift := range dailyShifts {
			totalHours += shift.HoursWorked
		}

		if totalHours == 0 {
			log.Printf("Total hours worked is zero for date: %s", date)
			continue
		}

		for _, shift := range dailyShifts {
			tips := (shift.HoursWorked / totalHours) * dailyTip.TotalTips
			tips = math.Round(tips*100) / 100 // Round to 2 decimal places
			log.Printf("WorkerID: %s, Date: %s, HoursWorked: %f, TotalHours: %f, TotalTips: %f, TipsEarned: %f",
				shift.WorkerID, shift.Date, shift.HoursWorked, totalHours, dailyTip.TotalTips, tips)
			tipsEarned = append(tipsEarned, &models.WorkerTip{
				WorkerID:    shift.WorkerID,
				LocationID:  shift.LocationID,
				Date:        shift.Date,
				HoursWorked: shift.HoursWorked,
				TipsEarned:  tips,
			})
		}
	}

	return tipsEarned, nil
}
