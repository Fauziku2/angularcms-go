package pages

import "github.com/gorilla/mux"

// PageHandler APIs
func PageHandler(r *mux.Router) {
	pages := r.PathPrefix("/pages").Subrouter()
	pages.HandleFunc("", getAllPages).Methods("GET")
	pages.HandleFunc("/{slug}", getPage).Methods("GET")
	pages.HandleFunc("/add-page", addPage).Methods("POST")
	pages.HandleFunc("/edit-page/{pageId}", getEditPage).Methods("GET")
	pages.HandleFunc("/edit-page/{pageId}", updatePage).Methods("PUT")
	pages.HandleFunc("/delete-page/{pageId}", deletePage).Methods("DELETE")
}
