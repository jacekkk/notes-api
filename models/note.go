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

func GetNotesByUser(uid int) ([]Note, error) {
	notesCollection := client.Database("notes").Collection("notes")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := notesCollection.Find(ctx, bson.M{"uid": uid})

	if err != nil {
		return nil, err
	}

	var userNotes []Note

	if err = cursor.All(ctx, &userNotes); err != nil {
		return nil, err
	}

	return userNotes, nil
}
