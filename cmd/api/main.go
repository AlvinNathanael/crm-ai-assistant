package main

import (
	"log"
	"net/http"

	"github.com/AlvinNathanael/crm-ai-assistant/internal/app"
)

func main() {
	app := app.NewApp()

	// Listen and serve.
	log.Println("Server starting on: 8080")
	if err := http.ListenAndServe(":8080", app.Mux); err != nil {
		log.Fatalln("Failed HTTP server: " + err.Error())
	}
}
