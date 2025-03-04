package main

import (
	"trinkgeldApp/database"
	"trinkgeldApp/handlers"
	"trinkgeldApp/server"
)

func main() {
	// Initialize PocketBase
	pb := database.InitializePocketBase()

	// Start echo

	// Create AppContext that holds database
	appContext := &handlers.AppContext{
		DB: pb,
	}

	// Initialize Echo Server and pass appContext
	e := server.InitializeEchoServer(appContext)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))

}
