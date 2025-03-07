package handlers

import (
	"log"
	"reflect"
	"trinkgeldApp/models"

	"github.com/pocketbase/pocketbase/core"
)

func (a *AppContext) UploadShiftsForPeriod(shifts []*models.WorkShift) {
	// Upload the shifts to the database
	collection, err := a.DB.FindCollectionByNameOrId("work_shifts")
	if err != nil {
		log.Fatalf("Error finding collection: %v", err)
	}

	for _, shift := range shifts {
		record := core.NewRecord(collection)
		val := reflect.ValueOf(shift).Elem()
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fieldName := field.Tag.Get("json")
			fieldValue := val.Field(i).Interface()
			record.Set(fieldName, fieldValue)
		}
		err := a.DB.Save(record)
		if err != nil {
			log.Fatalf("Error saving record: %v", err)
		}
	}
}

func (a *AppContext) UploadTipsForPeriod(tips []*models.DailyTip, location string) {
	// Upload the tips to the database
	collection, err := a.DB.FindCollectionByNameOrId("daily_tips")
	if err != nil {
		log.Fatalf("Error finding collection: %v", err)
	}

	for _, tip := range tips {
		log.Println(tip)
		record := core.NewRecord(collection)
		val := reflect.ValueOf(tip).Elem()
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fieldName := field.Tag.Get("json")
			fieldValue := val.Field(i).Interface()
			record.Set(fieldName, fieldValue)
		}
		// Set the location_id field with the location variable
		record.Set("location_id", location)
		err := a.DB.Save(record)
		if err != nil {
			log.Fatalf("Error saving record: %v", err)
		}
	}
}
