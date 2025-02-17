package main

import (
	"log"
	"news/db"
	"news/handlers"
	"news/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	utils.InitLogger()

	db.InitDB()

	app := fiber.New()

	app.Post("/edit/:Id", handlers.EditNews)
	app.Get("/list", handlers.ListNews)

	utils.Log.Info("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
