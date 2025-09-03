package src

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/vandimit/simple-hotels-mock-rest-api/src/handlers"
	"github.com/vandimit/simple-hotels-mock-rest-api/src/services"
	"github.com/vandimit/simple-hotels-mock-rest-api/src/utils"
)

func Main() {
	// Create services
	hotelService := services.NewHotelService()

	// Load hotel data from JSON file
	dataPath := "mock-data/hotels-data.json"
	err := hotelService.LoadHotelsFromFile(dataPath)
	if err != nil {
		log.Fatalf("Failed to load hotel data: %v", err)
	}

	// Create handlers
	hotelHandler := handlers.NewHotelHandler(hotelService)

	// Create router
	router := mux.NewRouter()
	
	// Add middleware
	router.Use(utils.LoggingMiddleware)
	router.Use(utils.CORSMiddleware)

	// API routes with prefix
	apiRouter := router.PathPrefix("/api").Subrouter()
	
	// Register routes
	apiRouter.HandleFunc("/hotels", hotelHandler.GetHotels).Methods("GET")
	apiRouter.HandleFunc("/hotels/{hotelId}", hotelHandler.GetHotelByID).Methods("GET")

	// Set up server
	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	// Set up graceful shutdown
	gracefulShutdown(srv)
}

// gracefulShutdown handles graceful server shutdown on interrupt signal
func gracefulShutdown(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until signal is received
	<-c

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Shutdown server
	fmt.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")
	os.Exit(0)
}