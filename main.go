package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gocroot/url"
    "github.com/gocroot/controller"
)

func main() {
	// Inisialisasi router
	router := mux.NewRouter()

	// Menambahkan rute
	router.HandleFunc("/", url.URL).Methods("GET")
	router.HandleFunc("/api/nearby-roads", controller.NearbyRoadHandler).Methods("GET")

	// Memulai server pada port 8080
	log.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
