package main

import (
	"fmt"
	"formality/backend/database"
	"formality/backend/listener"
	migration "formality/backend/migrations"
	"formality/backend/routes"
	"sync"

	"github.com/joho/godotenv"
)


func main() {
	migration.Migration()
	database.ConnectDatabase()
	err := godotenv.Load()//by default, it is .env so we don't have to write
   if err != nil {
      fmt.Println("Error is occurred  on .env file please check")
   }

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		routes.Routes()
	}()

	go func() {
		defer wg.Done()
		listener.Listener()
	}()

	wg.Wait()
}