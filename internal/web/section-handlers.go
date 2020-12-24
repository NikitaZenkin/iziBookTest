package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (c *Controller) CreateSubSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	parentSectionId := chi.URLParam(r, urlSectionId)
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

	newSectionID, err := c.repository.CreateSection(ctx, parentSectionId, data.Name)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, map[string]string{"section_id": newSectionID}, http.StatusOK)
}

func (c *Controller) CreateRootSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	newSectionID, err := c.repository.CreateSection(ctx, "", data.Name)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, map[string]string{"section_id": newSectionID}, http.StatusOK)
}

func (c *Controller) ReadSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sectionId := chi.URLParam(r, urlSectionId)
	if sectionId == "" {
		RespondWithError(w, fmt.Errorf("no url params"), http.StatusBadRequest)
		return
	}

	ok, err := c.repository.CheckSection(ctx, sectionId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		RespondWithError(w, fmt.Errorf("section is not yours or does not exist"), http.StatusForbidden)
		return
	}

	section, err := c.repository.GetSectionTree(ctx, sectionId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, section, http.StatusOK)
}

func (c *Controller) UpdateSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	sectionId := chi.URLParam(r, urlSectionId)
	if sectionId == "" {
		RespondWithError(w, fmt.Errorf("no url params"), http.StatusBadRequest)
		return
	}

	ok, err := c.repository.CheckSection(ctx, sectionId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		RespondWithError(w, fmt.Errorf("section is not yours or does not exist"), http.StatusForbidden)
		return
	}

	err = c.repository.UpdateSection(ctx, sectionId, data.Name)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteSection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sectionId := chi.URLParam(r, urlSectionId)
	if sectionId == "" {
		RespondWithError(w, fmt.Errorf("no url params"), http.StatusBadRequest)
		return
	}

	ok, err := c.repository.CheckSection(ctx, sectionId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		RespondWithError(w, fmt.Errorf("section is not yours or does not exist"), http.StatusForbidden)
		return
	}

	err = c.repository.DeleteSection(ctx, sectionId)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
