package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rifqoi/sns-go/internal/config"
	"github.com/rifqoi/sns-go/postgres"
	"github.com/rifqoi/sns-go/rest"
	"github.com/rifqoi/sns-go/services/user"
)

func main() {
	config.LoadEnv(".env")
	r := chi.NewRouter()

	db, err := postgres.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	userService, err := user.NewUserService(user.WithPostgresUserRepository(db))
	if err != nil {
		panic(err)
	}

	rest.NewUserController(userService).Register(r)
	s := http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("Starting server on port 3000")

		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server %s\n", err.Error())
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c

	log.Println("Got signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}
