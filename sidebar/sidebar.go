package sidebar

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Fauziku2/ngExpressCms/goapi/database"
	"github.com/Fauziku2/ngExpressCms/goapi/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Sidebar model
type Sidebar struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content" bson:"content"`
}

//Get a sidebar
func getSidebar(w http.ResponseWriter, r *http.Request) {
	id := handlers.MustGetObjectID(w, "5e4783384492bf62dc155f9f")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sidebar Sidebar
	err2 := database.DB.Sidebar.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&sidebar)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err2.Error() + `}`))
		log.Fatalln(err2)
	}

	handlers.Write(w, http.StatusOK, sidebar)
}

// Post add sidebar
func postSidebar(w http.ResponseWriter, r *http.Request) {
	var sideBar Sidebar
	json.NewDecoder(r.Body).Decode(&sideBar)

	id := handlers.MustGetObjectID(w, "5e4783384492bf62dc155f9f")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err2 := database.DB.Sidebar.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"content": sideBar.Content}})
	if err2 != nil {
		w.Write([]byte(`{"message: "` + err2.Error() + `}`))
		handlers.Write(w, http.StatusInternalServerError, "problem")
		log.Fatalln(err2)
	}
	handlers.Write(w, http.StatusOK, "ok")
}
