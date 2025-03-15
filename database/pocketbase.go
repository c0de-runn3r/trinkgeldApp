package database

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func InitializePocketBase() *pocketbase.PocketBase {

	app := pocketbase.New()

	app.OnServe().BindFunc(func(event *core.ServeEvent) error {
		// This will run when the server is ready to serve requests
		log.Println("PocketBase is ready to serve requests")
		return event.Next()
	})

	// Start PocketBase in the background
	go func() {
		if err := app.Start(); err != nil {
			log.Fatalf("Error starting PocketBase: %v", err)
		}
	}()

	return app
}
