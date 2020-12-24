package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"iziBookTest/internal/util"
)

func (c *Controller) Registration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &util.User{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	loginExist, err := c.repository.LoginExist(ctx, data.Login)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
	if loginExist {
		RespondWithError(w, fmt.Errorf("login already exist"), http.StatusForbidden)
		return
	}

	userID, err := c.repository.CreateUser(ctx, data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	err = c.sessionManager.CreateSession(ctx, w, userID)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) LogIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	loginExist, err := c.repository.LoginExist(ctx, data.Login)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
	if !loginExist {
		RespondWithError(w, fmt.Errorf("login does not exist exist"), http.StatusUnauthorized)
		return
	}

	userID, ok, err := c.repository.CheckPassWord(ctx, data.Login, data.Password)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}
	if !ok {
		RespondWithError(w, fmt.Errorf("wrong password"), http.StatusUnauthorized)
		return
	}

	err = c.sessionManager.CreateSession(r.Context(), w, userID)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) LogOut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &util.User{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	err = c.sessionManager.LogOut(ctx, w)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data := &util.User{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	err = c.repository.UpdateUser(ctx, data)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := c.sessionManager.LogOut(ctx, w)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	err = c.repository.DeleteUser(ctx)
	if err != nil {
		c.log.Println(err)
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
