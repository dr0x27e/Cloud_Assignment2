package webhooks

import (
	"Assignment2/internal/services"
	"net/http"
)

func Notifications(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		services.WebhookReg(w, r)
	case http.MethodDelete:
		services.DeleteWeb(w, r)
	case http.MethodGet:
		// Checking if ID is provided and decides which to use.
		if id := r.PathValue("id"); id != "" {
			services.GET_Id_Webhook(w, r)
		} else {
			services.GetAllWeb(w, r)
		}

	default:
		Service("REGISTER", "NO")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

