package gastromaticprocessor

import (
	"fmt"
	"strconv"
	"time"
	"trinkgeldApp/models"

	"github.com/xuri/excelize/v2"
)

const (
	GastromaticDateLayout = "02.01.2006" // to parse the date without the weekday
)

func ProcessGastromaticFile(name string) ([]*models.WorkShift, error) {
	file, err := openGastromaticFile(name)
	if err != nil {
		return nil, err
	}

	dailyWorkingTimesPerWorker, err := extractDailyWorkingTimesPerWorker(file)
	if err != nil {
		return nil, err
	}

	return dailyWorkingTimesPerWorker, nil
}

func openGastromaticFile(name string) (*excelize.File, error) {
	f, err := excelize.OpenFile(name)
	if err != nil {
		return nil, err
	}

	// jsut to be save that it was closed correctly
	defer func() error {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			return err
		}
		return nil
	}()

	sheetList := f.GetSheetList()
	if len(sheetList) == 1 {
		return nil, fmt.Errorf("File that you are trying to upload is not Gastromatic detailed export, it contains only one sheet")
	}

	if sheetList[0] != "Übersicht" {
		return nil, fmt.Errorf("File that you are trying to upload is not Gastromatic detailed export, the first sheet is not Übersicht")
	}

	return f, nil
}

func extractDailyWorkingTimesPerWorker(file *excelize.File) ([]*models.WorkShift, error) {
	workers := file.GetSheetList()[1:]

	var dailyWorkingTimesPerWorker []*models.WorkShift

	for _, worker := range workers {
		rows, err := file.GetRows(worker)
		if err != nil {
			return nil, err
		}

		// Find the column index for "Dauer netto (dezimal)"
		hoursWorkedColIndex := -1
		for colIndex, cellValue := range rows[5] {
			if cellValue == "Dauer netto (dezimal)" {
				hoursWorkedColIndex = colIndex
				break
			}
		}

		if hoursWorkedColIndex == -1 {
			return nil, fmt.Errorf("could not find 'Dauer netto (dezimal)' column in sheet %s", worker)
		}

		for i, row := range rows {
			if i < 6 { // Skip the first 6 rows (0-based index)
				continue
			}

			if len(row) > 0 && row[0] == "Summe:" {
				break
			}

			// Process the row
			if len(row) > 1 {
				dateString := row[0]
				workingType := row[1]  // working/sick/planned/absent etc. Working should be "A"
				location := row[4]     // location name
				positionType := row[5] // position type
				hoursWorkedString := row[hoursWorkedColIndex]

				if workingType == "A" && positionType == "Barista" { // check if the person was working and was on the barista position

					// Remove the short name of the workday (2 letters and a space)
					if len(dateString) > 3 {
						dateString = dateString[3:]
					}

					// Convert dateString to time.Time
					date, err := time.Parse(GastromaticDateLayout, dateString)
					if err != nil {
						return nil, fmt.Errorf("error parsing date at row %d: %w", i, err)
					}

					// Convert hoursWorked to float64
					hoursWorked, err := strconv.ParseFloat(hoursWorkedString, 64)
					if err != nil {
						return nil, fmt.Errorf("error converting hours worked at row %d: %w", i, err)
					}

					dailyWorkingTimesPerWorker = append(dailyWorkingTimesPerWorker, &models.WorkShift{
						WorkerID:    worker,   // I use the worker name as the id of the worker
						LocationID:  location, // I use the location name as the id of the location
						Date:        date,
						HoursWorked: hoursWorked,
					})
				}
			}
		}
	}

	return dailyWorkingTimesPerWorker, nil
}
