package sidebar

import "github.com/gorilla/mux"

// SidebarHandler APIs
func SidebarHandler(r *mux.Router) {
	sidebar := r.PathPrefix("/sidebar").Subrouter()
	sidebar.HandleFunc("/edit-sidebar", getSidebar).Methods("GET")
	sidebar.HandleFunc("/edit-sidebar", postSidebar).Methods("POST")
}
