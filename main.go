package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Fauziku2/ngExpressCms/goapi/database"
	"github.com/Fauziku2/ngExpressCms/goapi/pages"
	"github.com/Fauziku2/ngExpressCms/goapi/sidebar"
	"github.com/Fauziku2/ngExpressCms/goapi/users"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Init Router
	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// mongo database
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	database.NewDB(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	// Route Handlers / Endpoints
	users.UserHandler(r)
	pages.PageHandler(r)
	sidebar.SidebarHandler(r)

	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(r)))
}
