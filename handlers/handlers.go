package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MustGetObjectID converts string to primitive.ObjectID
func MustGetObjectID(w http.ResponseWriter, id string) primitive.ObjectID {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err.Error() + `}`))
		log.Fatalln(err)
	}

	return ID
}

// Write Handler
func Write(w http.ResponseWriter, code int, out interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(out)
	if err != nil {
		panic(err)
	}
}
