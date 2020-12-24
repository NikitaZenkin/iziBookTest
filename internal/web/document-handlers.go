package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"iziBookTest/internal/util"
)

func (c *Controller) CreateDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &util.Document{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	parentSectionId := chi.URLParam(r, urlDocumentId)
	if parentSectionId == "" {
		RespondWithError(w, fmt.Errorf("no url params"), http.StatusBadRequest)
		return
	}

	ok, err := c.repository.CheckSection(ctx, parentSectionId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		RespondWithError(w, fmt.Errorf("parent section is not yours or does not exist"), http.StatusForbidden)
		return
	}

	newSectionID, err := c.repository.CreateDocument(ctx, parentSectionId, data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, map[string]string{"document_id": newSectionID}, http.StatusOK)
}

func (c *Controller) ReadDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentId := chi.URLParam(r, urlDocumentId)
	if documentId == "" {
		RespondWithError(w, fmt.Errorf("no url params"), http.StatusBadRequest)
		return
	}

	ok, err := c.repository.CheckDocument(ctx, documentId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		RespondWithError(w, fmt.Errorf("document is not yours or does not exist"), http.StatusForbidden)
		return
	}

	document, err := c.repository.GetDocument(ctx, documentId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, document, http.StatusOK)
}

func (c *Controller) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &util.Document{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	documentId := chi.URLParam(r, urlDocumentId)
	if documentId == "" {
		RespondWithError(w, fmt.Errorf("no url params"), http.StatusBadRequest)
		return
	}

	ok, err := c.repository.CheckDocument(ctx, documentId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		RespondWithError(w, fmt.Errorf("document is not yours or does not exist"), http.StatusForbidden)
		return
	}

	err = c.repository.UpdateDocument(ctx, documentId, data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	documentId := chi.URLParam(r, urlDocumentId)
	if documentId == "" {
		RespondWithError(w, fmt.Errorf("no url params"), http.StatusBadRequest)
		return
	}

	ok, err := c.repository.CheckDocument(ctx, documentId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		RespondWithError(w, fmt.Errorf("document is not yours or does not exist"), http.StatusForbidden)
		return
	}

	err = c.repository.DeleteDocument(ctx, documentId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
