package main

import (
	"log"
	"net/http"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/app"
	"github.com/AlvinNathanael/crm-ai-assistant/internal/config"
)

func main() {
	// Load config.
	config := config.Load()

	// Init app.
	app, err := app.NewApp(config)
	if err != nil {
		log.Fatalln("Failed to initialize app: " + err.Error())
	}
	defer app.DB.Close()

	// Listen and serve.
	log.Println("Server starting on: 8080")
	if err := http.ListenAndServe(":8080", app.Mux); err != nil {
		log.Fatalln("Failed HTTP server: " + err.Error())
	}
}
