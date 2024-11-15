package handlers

import "net/http"

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
		"Message": "unexpected error",
	})
}
