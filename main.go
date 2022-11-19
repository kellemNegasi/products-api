package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/kellemNegasi/product-api/handlers"
)

const address = ":8080"

func main() {
	l := log.New(os.Stdout, "Products-api ", log.LstdFlags)
	ph := handlers.NewProducts(l)
	r := mux.NewRouter()

	getRouter := r.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	postRouter := r.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)
	putRouter := r.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	server := &http.Server{
		Addr:         address,
		Handler:      r,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  2 * time.Second,
	}

	go func() {
		fmt.Println("starting server at ", address)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	s := <-sigChan

	fmt.Println("\nReceived terminate signal shutting down gracefully ", s)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
