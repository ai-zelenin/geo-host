package server

import "net/http"

func RedirectPermanent(newUrl string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, newUrl, http.StatusMovedPermanently)
	})
}
