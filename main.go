package main

import (
	"formality/migrations"
	"formality/packages/database"
	"formality/packages/listener"
	"formality/packages/routes"
	"sync"
)

func main() {
	migrations.Migration()
	database.ConnectDatabase()

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
