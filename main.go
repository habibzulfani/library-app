package main

import (
	"os"
	"project/configs"
	"project/routes"
)

func main() {

	configs.InitDB()

	e := routes.New()
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Default port if not specified in .env
	}
	e.Logger.Fatal(e.Start(":" + port))
}
