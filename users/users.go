package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Fauziku2/ngExpressCms/goapi/database"
	"github.com/Fauziku2/ngExpressCms/goapi/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User model
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
}

func test(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("this get call all the time")
		next.ServeHTTP(w, r)
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	var body User
	json.NewDecoder(r.Body).Decode(&body)

	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	var user *User
	err1 := database.DB.Users.FindOne(ctx1, bson.M{
		"username": body.Username,
	}).Decode(&user)
	if err1 != nil && err1 != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err1.Error() + `}`))
		panic(err1)
	}
	if user != nil {
		handlers.Write(w, http.StatusOK, "userExist")
		return
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	_, err2 := database.DB.Users.InsertOne(ctx2, body)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err2.Error() + `}`))
		panic(err2)
	}

	handlers.Write(w, http.StatusOK, "userRegistered")
}

func login(w http.ResponseWriter, r *http.Request) {
	var body User
	json.NewDecoder(r.Body).Decode(&body)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user *User
	err := database.DB.Users.FindOne(ctx, bson.M{
		"username": body.Username,
		"password": body.Password,
	}).Decode(&user)

	if err != nil && err != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err.Error() + `}`))
		panic(err)
	}

	if user != nil {
		handlers.Write(w, http.StatusOK, body.Username)
		return
	}

	handlers.Write(w, http.StatusOK, "invalidLogin")
}
