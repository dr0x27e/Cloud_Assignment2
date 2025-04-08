package api

import (
	"Assignment2/internal/services"
	"net/http"
	"log"
)

func Registrations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		services.POST_Registration(w, r)
	case http.MethodGet:
		// Checking if ID is provided and decides which to use.
		if id := r.PathValue("id"); id != "" {
			services.GET_Id_Registration(w, r)
		} else {
			services.GET_All_Registration(w, r)
		}
	case http.MethodDelete:
		services.DELETE_Registration(w, r)
	case http.MethodPut:
		services.PUT_Registration(w, r)
	case http.MethodPatch:
		services.PATCH_Registration(w, r)
	default:
		log.Println("Unsupported request method " + r.Method)
		http.Error(w,
			"Unsupported request method "+r.Method,
			http.StatusMethodNotAllowed)
		return
	}
}
