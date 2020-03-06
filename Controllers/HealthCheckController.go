package Controllers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type HealthCheckController struct {
}

func NewHealthCheckController() *HealthCheckController {
	controller := HealthCheckController{}

	return &controller
}

func (this *HealthCheckController) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", this.Get)
	return router
}

func (this *HealthCheckController) Get(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, true)
}
