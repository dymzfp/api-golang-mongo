package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	ctrl "golang-mongo/controller"
)

const (
	PORT = ":8081"
)

func Handler() {
	r := mux.NewRouter()

	r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Home")
	})
	r.HandleFunc("/api", ctrl.PostData).Methods(http.MethodPost)
	r.HandleFunc("/api", ctrl.GetData).Methods(http.MethodGet)
	r.HandleFunc("/api/{id}", ctrl.GetDataSingle).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(PORT, r))
}

func main() {
	log.Printf("Server running at port %v", PORT)
	Handler()
}