package gastronoviprocessor_test

import (
	"fmt"
	"testing"
	gastronoviprocessor "trinkgeldApp/gastronoviProcessor"
)

func TestProcessGastronoviFile(t *testing.T) {
	t.Run("Test ProcessGastronoviFile", func(t *testing.T) {
		tips, err := gastronoviprocessor.ProcessGastronoviFile("testdata/test.xlsx")
		if err != nil {
			t.Errorf("Error while processing the file: %v", err)
		} else {
			for _, tip := range tips {
				fmt.Printf("Date: %s, TotalTips: %.2f\n", tip.Date, tip.TotalTips)
			}
		}
		if len(tips) == 0 {
			t.Errorf("No tips were extracted")
		}
	})
}
