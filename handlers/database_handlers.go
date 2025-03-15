package handlers

import (
	"encoding/json"
	"log"
	mainprocessor "trinkgeldApp/mainprocessor"
	"trinkgeldApp/models"
	"trinkgeldApp/utils"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func (a *AppContext) UploadShiftsForPeriod(shifts []*models.WorkShift) {
	// Upload the shifts to the database
	collection, err := a.DB.FindCollectionByNameOrId(models.WorkShiftCollection)
	if err != nil {
		log.Fatalf("Error finding collection: %v", err)
	}

	// TODO can use marshalling instead of reflection
	for _, shift := range shifts {

		// Check if the worker already exists in the database
		workerExists, err := a.getOrCreateWorker(shift.WorkerID)
		if err != nil {
			log.Fatalf("Error getting or creating worker: %v", err)
		}
		if !workerExists {
			log.Fatalf("Error getting or creating worker: %v", err)
		}

		// before we receive worker name, we need to generate worker ID and assign it to the shift
		shift.WorkerID = utils.GenerateWorkerID(shift.WorkerID)

		exsistingShift, err := a.DB.FindAllRecords(models.WorkShiftCollection, dbx.HashExp{"worker_id": shift.WorkerID, "date": shift.Date, "location_id": shift.LocationID})
		if err != nil && err.Error() != "sql: no rows in result set" {
			log.Fatalf("Error during the search looking into work_shifts: %v", err)
		}

		switch len(exsistingShift) {
		case 0:
			// do nothing
		case 1:
			log.Printf("Shift already exists, updating the record")
			shift.ID = exsistingShift[0].Id
		default:
			log.Fatalf("Error: multiple records found for worker_id: %s, date: %s", shift.WorkerID, shift.Date)
		}

		record := core.NewRecord(collection)
		record.Set("id", shift.ID)
		record.Set("worker_id", shift.WorkerID)
		record.Set("location_id", shift.LocationID)
		record.Set("date", shift.Date)
		record.Set("hours_worked", shift.HoursWorked)

		err = a.DB.Save(record)
		if err != nil {
			log.Fatalf("Error saving record: %v", err)
		}
	}

}

func (a *AppContext) UploadTipsForPeriod(tips []*models.DailyTip, location string) {
	// Upload the tips to the database
	collection, err := a.DB.FindCollectionByNameOrId(models.DailyTipCollection)
	if err != nil {
		log.Fatalf("Error finding collection: %v", err)
	}

	// TODO can use marshalling instead of reflection
	for _, tip := range tips {

		tipExists, err := a.DB.FindAllRecords(models.DailyTipCollection, dbx.HashExp{"date": tip.Date, "location_id": location})
		if err != nil && err.Error() != "sql: no rows in result set" {
			log.Fatalf("Error during the search looking into daily_tips: %v", err)
		}

		switch len(tipExists) {
		case 0:
			// do nothing
		case 1:
			log.Printf("Tip already exists, updating the record")
			tip.ID = tipExists[0].Id
		default:
			log.Fatalf("Error: multiple records found for date: %s", tip.Date)

			record := core.NewRecord(collection)
			record.Set("id", tip.ID)
			record.Set("date", tip.Date)
			record.Set("total_tips", tip.TotalTips)
			record.Set("location_id", location)
			err := a.DB.Save(record)
			if err != nil {
				log.Fatalf("Error saving record: %v", err)
			}
		}
	}
}

func (a *AppContext) getOrCreateWorker(workerName string) (bool, error) {
	workerID := utils.GenerateWorkerID(workerName)

	// Check if the worker already exists in the database
	collection, err := a.DB.FindCollectionByNameOrId(models.WorkerCollection)
	if err != nil {
		log.Printf("Error finding collection: %v", err)
		return false, err
	}

	result, err := a.DB.FindRecordById("workers", workerID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Printf("Worker not found, creating new worker with ID: %s", workerID)
		} else {
			log.Printf("Error finding worker: %v", err)
			return false, err
		}
	} else if result != nil && result.Id == workerID {
		log.Printf("Worker found with ID: %s", workerID)
		return true, nil
	}

	// Create a new worker record if it does not exist
	record := core.NewRecord(collection)
	record.Set("id", workerID)
	record.Set("name", workerName)
	err = a.DB.Save(record)
	if err != nil {
		log.Printf("Error saving new worker: %v", err)
		return false, err
	}
	log.Printf("New worker created with ID: %s", workerID)
	return true, nil
}

func (a *AppContext) GetWorkShifts() ([]*models.WorkShift, error) {
	collection, err := a.DB.FindCollectionByNameOrId(models.WorkShiftCollection)
	if err != nil {
		return nil, err
	}

	records, err := a.DB.FindAllRecords(collection)
	if err != nil {
		return nil, err
	}

	var workShifts []*models.WorkShift
	for _, record := range records {
		var shift models.WorkShift
		data, err := record.MarshalJSON()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &shift)
		if err != nil {
			return nil, err
		}
		workShifts = append(workShifts, &shift)
	}

	return workShifts, nil
}

func (a *AppContext) GetDailyTips() ([]*models.DailyTip, error) {
	collection, err := a.DB.FindCollectionByNameOrId(models.DailyTipCollection)
	if err != nil {
		return nil, err
	}

	records, err := a.DB.FindAllRecords(collection)
	if err != nil {
		return nil, err
	}

	var dailyTips []*models.DailyTip
	for _, record := range records {
		var tip models.DailyTip
		data, err := record.MarshalJSON()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &tip)
		if err != nil {
			return nil, err
		}
		dailyTips = append(dailyTips, &tip)
	}
	return dailyTips, nil
}

func (a *AppContext) UploadTipsForPeriodPerWorkerPerDay(tipsEarned []*models.WorkerTip) {

	collection, err := a.DB.FindCollectionByNameOrId(models.WorkerTipCollection)
	if err != nil {
		log.Fatalf("Error finding collection: %v", err)
	}

	for _, tip := range tipsEarned {

		record := core.NewRecord(collection)

		// Check if a record with the same worker_id and date already exists
		existingRecords, err := a.DB.FindAllRecords(models.WorkerTipCollection, dbx.HashExp{"worker_id": tip.WorkerID, "date": tip.Date, "location_id": tip.LocationID})
		if err != nil && err.Error() != "sql: no rows in result set" {
			log.Fatalf("Error during the search looking into worker_tips: %v", err)
		}

		switch len(existingRecords) {
		case 0:
			record = nil
		case 1:
			record = existingRecords[0]
		default:
			log.Fatalf("Error: multiple records found for worker_id: %s, date: %s", tip.WorkerID, tip.Date)
		}

		record.Set("worker_id", tip.WorkerID)
		record.Set("date", tip.Date)
		record.Set("tips_earned", tip.TipsEarned)
		record.Set("hours_worked", tip.HoursWorked)
		record.Set("location_id", tip.LocationID)

		err = a.DB.Save(record)
		if err != nil {
			log.Fatalf("Error saving record: %v", err)
		}

	}

}

func (a *AppContext) CheckDBandCalculateTips() {
	// Get all the work shifts from the database
	workShifts, err := a.GetWorkShifts()
	if err != nil {
		log.Fatalf("Error getting work shifts: %v", err)
	}

	// Get all the daily tips from the database
	dailyTips, err := a.GetDailyTips()
	if err != nil {
		log.Fatalf("Error getting daily tips: %v", err)
	}

	tipsEarned, err := mainprocessor.CalculateTipAmountsPerWorkerPerDay(workShifts, dailyTips)
	if err != nil {
		log.Fatalf("Error calculating tips: %v", err)
	}

	// Upload the calculated tips to the database
	a.UploadTipsForPeriodPerWorkerPerDay(tipsEarned)
}
