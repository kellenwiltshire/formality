package main

import (
	"formality/backend/database"
	migration "formality/backend/migrations"
	"formality/backend/routes"
)


func main() {
	migration.Migration()
	database.ConnectDatabase()
	routes.Routes()
}