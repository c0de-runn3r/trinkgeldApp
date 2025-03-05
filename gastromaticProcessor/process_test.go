package gastromaticprocessor_test

import (
	"fmt"
	"testing"
	gastromaticprocessor "trinkgeldApp/gastromaticProcessor"
)

func TestProcessGastromaticFile(t *testing.T) {
	t.Run("Test ProcessGastromaticFile", func(t *testing.T) {
		workingTimes, err := gastromaticprocessor.ProcessGastromaticFile("testdata/test.xlsx")
		if err != nil {
			t.Errorf("Error while processing the file: %v", err)
		} else {
			for _, workingTime := range workingTimes {
				fmt.Printf("Worker: %s, Date: %s, WorkingTime: %.2f, Location: %s\n", workingTime.WorkerID, workingTime.Date.Format("02.01.2006"), workingTime.HoursWorked, workingTime.LocationID)
			}
		}
		if len(workingTimes) == 0 {
			t.Errorf("No working times were extracted")
		}
	})
}
