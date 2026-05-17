package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to websocket", err)
		return
	}
	defer conn.Close()

	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		log.Println("User ID is required", http.StatusBadRequest)
		return
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message", err)
			break
		}
		log.Println("Received message", string(msg))
	}

}

func handleDriversWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to websocket", err)
		return
	}
	defer conn.Close()

	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		log.Println("User ID is required", http.StatusBadRequest)
		return
	}

	packageSlug := r.URL.Query().Get("package_slug")
	if packageSlug == "" {
		log.Println("Package Slug is required", http.StatusBadRequest)
		return
	}

	type Driver struct {
		Id             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profile_picture"`
		CarPlate       string `json:"car_plate"`
		PackageSlug    string `json:"package_slug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			Id:             userId,
			Name:           "jfj",
			ProfilePicture: "sdfa",
			CarPlate:       "dsfa45234f",
			PackageSlug:    packageSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Println("Error writing message", err)
		return
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message", err)
			break
		}
		log.Println("Received message", string(msg))
	}
}
