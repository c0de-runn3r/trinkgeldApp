package gastronoviprocessor

import (
	"fmt"
	"strconv"
	"time"
	"trinkgeldApp/models"
	"trinkgeldApp/utils"

	"github.com/xuri/excelize/v2"
)

const (
	Sheet                = "Worksheet" // predefined by Gastronovi export
	GastronoviDateLayout = "02.01."    // to parse the date
)

func ProcessGastronoviFile(name string) ([]*models.DailyTip, error) {
	file, err := openGastronoviFile(name)
	if err != nil {
		return nil, err
	}

	tips, err := extractTipAmounts(file)
	if err != nil {
		return nil, err
	}

	return tips, nil
}

func openGastronoviFile(name string) (*excelize.File, error) {
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

	cell, err := f.GetCellValue(Sheet, "A6")
	if err != nil {
		return nil, err
	}

	// Check if it's the correct one
	if cell != "Trinkgeld" {
		return nil, fmt.Errorf("File that you are trying to upload is not correct or missing the Tips data (Maybe not in German?)")
	}

	return f, nil
}

func extractTipAmounts(file *excelize.File) ([]*models.DailyTip, error) {

	// slice that holds all tips per day that are in the file
	tipForPeriodPerDay := *new([]*models.DailyTip)

	// Remove column with names of the raws, we not gonna need it, since we proved that it's the correct one already
	err := file.RemoveCol(Sheet, "A")
	if err != nil {
		return nil, err
	}

	raws, err := file.GetRows(Sheet)
	if err != nil {
		return nil, err
	}

	// this number gonna show us the amount of days + 1 (first column with total numbers).
	// ATTENTION: we just removed one column with the names. Check the code above.
	lengthOfRaw := len(raws[0])

	// we gonna start from 1, to skip the total amount (first column).
	for col := 1; col < lengthOfRaw; col++ {
		dayString := raws[0][col]
		tipString := raws[5][col]

		// Set the year to the current year
		currentYear := strconv.Itoa(time.Now().Year())
		day := dayString + currentYear

		var tip float64
		if tipString == "" {
			tip = 0
		} else {
			tip, err = utils.ConvertCurrencyToNumber(tipString)
			if err != nil {
				return nil, fmt.Errorf("error converting tip at column %d: %w", col, err)
			}
		}

		tipForPeriodPerDay = append(tipForPeriodPerDay, &models.DailyTip{
			Date:      day,
			TotalTips: tip,
		})
	}

	return tipForPeriodPerDay, nil
}
