package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	tripHttp "ride-sharing/services/trip-service/infrastructure/http"
	"ride-sharing/services/trip-service/infrastructure/repository"
	"ride-sharing/services/trip-service/service"
	"syscall"
	"time"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)

	// fare := &domain.RideFareModal{
	// 	UserID: "1",
	// }
	// t, err := svc.CreateTrip(context.Background(), fare)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(t)

	mux := http.NewServeMux()
	httpHandler := &tripHttp.HttpHandler{Service: svc}
	mux.HandleFunc("POST /trip/preview", httpHandler.HandlePreviewTrip)

	server := &http.Server{
		Addr:    ":8082", // Trip service port
		Handler: mux,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Trip Service listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server: %v", err)

	case sig := <-shutdown:
		log.Printf("Server is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the server gracefully: %v", err)
			server.Close()
		}
		log.Println("Server stopped")
	}
}
