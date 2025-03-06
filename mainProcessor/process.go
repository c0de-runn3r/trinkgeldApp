package mainprocessor

import "trinkgeldApp/models"

// This function should calculate the tips earned by each worker per day
// based on the hours worked and the total tips of the location
// It should return a slice of WorkerTip structs
// Each WorkerTip struct should contain the worker ID, location ID, date, hours worked, tips earned, and timestamps
// The tips earned should be calculated as follows:
// tips earned = (hours worked / total hours worked) * total tips
// where total hours worked is the sum of hours worked by all workers at the location on that day
// and total tips is the total tips of the location on that day
func CalculateTipAmountsPerWorkerPerDay(shifts []*models.WorkShift, dailyTip *models.DailyTip) ([]*models.WorkerTip, error) {
	tipsEarned := make([]*models.WorkerTip, 0)

	totalHours := 0.0
	for _, shift := range shifts {
		totalHours += shift.HoursWorked
	}

	for _, shift := range shifts {
		tipsEarned = append(tipsEarned, &models.WorkerTip{
			WorkerID:    shift.WorkerID,
			LocationID:  shift.LocationID,
			Date:        shift.Date,
			HoursWorked: shift.HoursWorked,
			TipsEarned:  (shift.HoursWorked / totalHours) * dailyTip.TotalTips,
		})
	}

	return tipsEarned, nil

}
