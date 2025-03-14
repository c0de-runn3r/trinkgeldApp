package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertCurrencyToNumber(input string) (float64, error) {
	// Remove the currency symbol and whitespace
	cleaned := strings.ReplaceAll(input, "€", "")
	cleaned = strings.TrimSpace(cleaned)

	// Replace comma with dot for float conversion
	cleaned = strings.ReplaceAll(cleaned, ",", ".")

	// Convert to float
	value, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("Error converting string to float:", err)
	}

	return value, nil
}

func GenerateWorkerID(workerName string) string {
	// Remove whitespace
	workerName = strings.TrimSpace(workerName)

	// Convert to lowercase
	workerName = strings.ToLower(workerName)

	// Replace whitespace with underscore
	workerName = strings.ReplaceAll(workerName, " ", "_")

	return workerName
}
