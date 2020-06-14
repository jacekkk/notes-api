package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notes/models"
	"strconv"

	"github.com/gorilla/mux"
)

type NotesResponse struct {
	Count int           `json:"count" binding:"required"`
	Notes []models.Note `json:"notes" binding:"required"`
}

// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}", getNotes).Methods("GET", "OPTIONS")

	return router
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Notes API")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, _ := mux.Vars(r)["id"]

	parsedID, parsingErr := strconv.Atoi(id)
	if parsingErr != nil {
		fmt.Println(parsingErr)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	notes, err := models.GetNotesByUser(parsedID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	response := NotesResponse{
		Count: len(notes),
		Notes: notes,
	}

	json.NewEncoder(w).Encode(response)
}
