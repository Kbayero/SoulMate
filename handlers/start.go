package handlers

import "net/http"

func StartHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "start.html", nil)
}
