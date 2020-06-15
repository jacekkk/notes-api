package router

import (
	"encoding/json"
	"log"
	"net/http"
	"notes/models"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotesResponse struct {
	Count int            `json:"count" binding:"required"`
	Notes []*models.Note `json:"notes" binding:"required"`
}

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}/notes", getUserNotes).Methods("GET", "OPTIONS")
	router.HandleFunc("/notes/{id}", getNote).Methods("GET", "OPTIONS")

	return router
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Notes API")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func getUserNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, _ := mux.Vars(r)["id"]

	parsedID, parsingErr := strconv.Atoi(id)
	if parsingErr != nil {
		log.Print(parsingErr)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	notes, err := models.GetNotesByUser(parsedID)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	response := NotesResponse{
		Count: len(notes),
		Notes: notes,
	}

	json.NewEncoder(w).Encode(response)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, _ := mux.Vars(r)["id"]

	parsedID, parsingErr := convertToObjectID(id)
	if parsingErr != nil {
		log.Print(parsingErr)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	note, err := models.GetNote(parsedID)
	if err != nil {
		log.Print(err)

		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
		} else {
			http.Error(w, http.StatusText(500), 500)
		}

		return
	}

	json.NewEncoder(w).Encode(note)
}

func convertToObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return objectID, nil
}
