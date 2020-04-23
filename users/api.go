package users

import "github.com/gorilla/mux"

// UserHandler APIs
func UserHandler(r *mux.Router) {
	r.Use(test)
	users := r.PathPrefix("/users").Subrouter()
	users.HandleFunc("/register", register).Methods("POST")
	users.HandleFunc("/login", login).Methods("POST")
}
