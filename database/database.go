package database

import "go.mongodb.org/mongo-driver/mongo"

// client variable
// var client *mongo.Client

var db *mongo.Database

// DB database
var DB Collections

// Collections struct
type Collections struct {
	Users   *mongo.Collection
	Pages   *mongo.Collection
	Sidebar *mongo.Collection
}

// NewDB function
func NewDB(client *mongo.Client) {
	db := client.Database("angularcms")
	DB = Collections{
		Users:   db.Collection("users"),
		Pages:   db.Collection("pages"),
		Sidebar: db.Collection("sidebar"),
	}
}
