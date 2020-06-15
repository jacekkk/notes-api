package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID     primitive.ObjectID `json:"noteId,omitempty" bson:"_id,omitempty" binding:"required"`
	UserID int                `json:"userId,omitempty" bson:"uid,omitempty"`
	Text   string             `json:"text,omitempty" bson:"text,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
}

func GetNotesByUser(uid int) ([]*Note, error) {
	var userNotes []*Note
	notesCollection := client.Database("notes").Collection("notes")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// TODO add error handling for when the uid doesn't exist (no matches found)
	cursor, err := notesCollection.Find(ctx, bson.M{"uid": uid})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &userNotes)
	if err != nil {
		return nil, err
	}

	return userNotes, nil
}

func GetNote(id primitive.ObjectID) (*Note, error) {
	var note *Note
	notesCollection := client.Database("notes").Collection("notes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := notesCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&note)
	if err != nil {
		return nil, err
	}

	return note, nil
}
