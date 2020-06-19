package router

import (
	"encoding/json"
	"log"
	"net/http"
	"notes/models"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotesResponse struct {
	Count int            `json:"count" binding:"required"`
	Notes []*models.Note `json:"notes" binding:"required"`
}

func Router() http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"Authorization"},
		Debug:          true,
	})

	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{id}/notes", getUserNotes).Methods("GET", "OPTIONS")
	router.HandleFunc("/notes/{id}", getNote).Methods("GET", "OPTIONS")

	handler := c.Handler(router)

	return handler
}

func homeHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func getUserNotes(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	id, _ := mux.Vars(req)["id"]

	parsedID, parsingErr := strconv.Atoi(id)
	if parsingErr != nil {
		log.Print(parsingErr)
		http.Error(res, http.StatusText(500), 500)
		return
	}

	notes, err := models.GetNotesByUser(parsedID)
	if err != nil {
		log.Print(err)
		http.Error(res, http.StatusText(500), 500)
		return
	}

	response := NotesResponse{
		Count: len(notes),
		Notes: notes,
	}

	json.NewEncoder(res).Encode(response)
}

func getNote(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	id, _ := mux.Vars(req)["id"]

	parsedID, parsingErr := convertToObjectID(id)
	if parsingErr != nil {
		log.Print(parsingErr)
		http.Error(res, http.StatusText(500), 500)
		return
	}

	note, err := models.GetNote(parsedID)
	if err != nil {
		log.Print(err)

		if err == mongo.ErrNoDocuments {
			http.NotFound(res, req)
		} else {
			http.Error(res, http.StatusText(500), 500)
		}

		return
	}

	json.NewEncoder(res).Encode(note)
}

func convertToObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return objectID, nil
}
