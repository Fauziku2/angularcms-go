package pages

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Fauziku2/ngExpressCms/goapi/database"
	"github.com/Fauziku2/ngExpressCms/goapi/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Page model
type Page struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title   string             `json:"title" bson:"title"`
	Slug    string             `json:"slug" bson:"slug"`
	Content string             `json:"content" bson:"content"`
	Sidebar string             `json:"sidebar" bson:"sidebar"`
}

// Body - Request body
type Body struct {
	Title      string `json:"title" bson:"title"`
	Content    string `json:"content" bson:"content"`
	HasSidebar bool   `json:"hasSidebar" bson:"hasSidebar"`
}

//Get all pages
func getAllPages(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pages []Page
	cursor, err := database.DB.Pages.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		panic(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var page Page
		cursor.Decode(&page)
		pages = append(pages, page)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		panic(err)
	}

	handlers.Write(w, http.StatusOK, pages)
}

//Get a page
func getPage(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var page Page
	err := database.DB.Pages.FindOne(ctx, bson.M{
		"slug": slug,
	}).Decode(&page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err.Error() + `}`))
		panic(err)
	}

	handlers.Write(w, http.StatusOK, page)
}

// Post add page
func addPage(w http.ResponseWriter, r *http.Request) {
	var body Body
	json.NewDecoder(r.Body).Decode(&body)

	var newPage Page
	newPage.Title = body.Title
	newPage.Slug = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(strings.ToLower(body.Title)), `-`)
	newPage.Content = body.Content
	if body.HasSidebar {
		newPage.Sidebar = "yes"
	} else {
		newPage.Sidebar = "no"
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	var page *Page
	err1 := database.DB.Pages.FindOne(ctx1, bson.M{
		"slug": newPage.Slug,
	}).Decode(&page)
	if err1 != nil && err1 != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err1.Error() + `}`))
		log.Fatalln(err1)
	}

	if page != nil {
		handlers.Write(w, http.StatusOK, "pageExist")
		return
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	_, err2 := database.DB.Pages.InsertOne(ctx2, newPage)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err2.Error() + `}`))
		log.Fatalln(err2)
	}

	handlers.Write(w, http.StatusOK, "ok")
}

// GET edit page
func getEditPage(w http.ResponseWriter, r *http.Request) {
	pageID := handlers.MustGetObjectID(w, mux.Vars(r)["pageId"])

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var page Page
	err2 := database.DB.Pages.FindOne(ctx, bson.M{
		"_id": pageID,
	}).Decode(&page)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err2.Error() + `}`))
		log.Fatalln(err2)
	}

	handlers.Write(w, http.StatusOK, page)
}

// Post edit page
func updatePage(w http.ResponseWriter, r *http.Request) {
	pageID := handlers.MustGetObjectID(w, mux.Vars(r)["pageId"])

	var body Body
	json.NewDecoder(r.Body).Decode(&body)

	var updatedPage Page
	updatedPage.Title = body.Title
	updatedPage.Slug = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(strings.ToLower(body.Title)), `-`)
	updatedPage.Content = body.Content
	if body.HasSidebar {
		updatedPage.Sidebar = "yes"
	} else {
		updatedPage.Sidebar = "no"
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel1()

	var page *Page
	err1 := database.DB.Pages.FindOne(ctx1, bson.M{
		"slug": updatedPage.Slug,
		"_id": bson.M{
			"$ne": pageID,
		},
	}).Decode(&page)
	if err1 != nil && err1 != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err1.Error() + `}`))
		log.Fatalln(err1)
	}
	if page != nil {
		handlers.Write(w, http.StatusOK, "pageExist")
		return
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	_, err2 := database.DB.Pages.UpdateOne(ctx2, bson.M{"_id": pageID}, bson.M{"$set": bson.M{
		"title":   updatedPage.Title,
		"slug":    regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(strings.ToLower(updatedPage.Title)), `-`),
		"content": updatedPage.Content,
		"sidebar": updatedPage.Sidebar,
	}})
	if err2 != nil {
		w.Write([]byte(`{"message: "` + err2.Error() + `}`))
		handlers.Write(w, http.StatusInternalServerError, "problem")
		log.Fatalln(err2)
	}

	handlers.Write(w, http.StatusOK, "ok")
}

// GET delete page
func deletePage(w http.ResponseWriter, r *http.Request) {
	pageID := handlers.MustGetObjectID(w, mux.Vars(r)["pageId"])

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err2 := database.DB.Pages.DeleteOne(ctx, bson.M{
		"_id": pageID,
	})
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message: "` + err2.Error() + `}`))
		log.Fatalln(err2)
	}

	handlers.Write(w, http.StatusOK, "ok")
}
