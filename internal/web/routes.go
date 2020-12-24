package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

const (
	urlSectionId  = "sectionId"
	urlDocumentId = "documentId"
)

func (c *Controller) Mount(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Post("/registration", c.Registration)
		r.Post("/login", c.LogIn)
	})
	r.Group(func(r chi.Router) {
		r.Use(c.authMiddleware)
		r.Get("/logout", c.LogOut)
		r.Route("/profile", func(r chi.Router) {
			r.Patch("/", c.UpdateUser)
			r.Delete("/", c.DeleteUser)
		})
		r.Route("/section", func(r chi.Router) {
			r.Route("/{"+urlSectionId+"}", func(r chi.Router) {
				r.Post("/", c.CreateSubSection)
				r.Get("/", c.ReadSection)
				r.Patch("/", c.UpdateSection)
				r.Delete("/", c.DeleteSection)
			})
			r.Post("/", c.CreateRootSection)
		})
		r.Route("/document", func(r chi.Router) {
			r.Route("/{"+urlDocumentId+"}", func(r chi.Router) {
				r.Post("/", c.CreateDocument)
				r.Get("/", c.ReadDocument)
				r.Patch("/", c.UpdateDocument)
				r.Delete("/", c.DeleteDocument)
			})
		})
	})
}

func RespondWithJSON(w http.ResponseWriter, data interface{}, code int) {
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, err error, code int) {
	RespondWithJSON(w, map[string]string{"error": err.Error()}, code)
}
