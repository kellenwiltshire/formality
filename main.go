package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kellenwiltshire/formality/internal/app"
	"github.com/kellenwiltshire/formality/internal/listener"
	"github.com/kellenwiltshire/formality/internal/routes"
)

func main() {

	// Todo: This will be replaced with using env variable for port in the future
	var port int
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.Db.Close()

	err = app.UserHandler.HandleCreateAdminUser()
	if err != nil {
		app.Logger.Printf("Error: CreateAdminUser %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		r := routes.Routes(app)
		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      r,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		}

		app.Logger.Printf("Application running on port %d\n", port)

		err = server.ListenAndServe()
		if err != nil {
			app.Logger.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := listener.Listener(app)
		if err != nil {
			app.Logger.Fatal(err)
		}
	}()

	wg.Wait()
}
