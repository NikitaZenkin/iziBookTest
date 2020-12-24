package web

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"iziBookTest/internal/config"
	"iziBookTest/internal/repository"
)

type Controller struct {
	log            *log.Logger
	repository     repository.Repository
	sessionManager repository.SessionManager
	router         chi.Router
}

func NewController(config *config.Config) (*Controller, error) {
	rep, err := repository.NewStorage(&config.DB)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()

	controller := &Controller{
		log:            log.New(os.Stdout, "logger: ", 0),
		repository:     rep,
		sessionManager: rep,
		router:         router,
	}

	controller.Mount(controller.router)

	return controller, nil
}

func (c *Controller) StartService() {
	_ = http.ListenAndServe(":80", c.router)
}
