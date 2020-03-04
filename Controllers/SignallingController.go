package Controllers

import (
	"github.com/alexandrevicenzi/go-sse"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strconv"
)

type SignallingController struct {
	sseServer *sse.Server
}

func NewSignallingController(sseController *sse.Server) *SignallingController {
	controller := SignallingController{
		sseServer: sseController,
	}

	return &controller
}

func (this *SignallingController) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{signal}", this.Get)
	return router
}

func (this *SignallingController) Get(w http.ResponseWriter, r *http.Request) {
	signalUnparsed := chi.URLParam(r, "signal")
	signal, err := strconv.ParseBool(signalUnparsed)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error Processing Request", http.StatusBadRequest)
		return
	}

	if signal {
		this.sseServer.SendMessage("/api/event", sse.SimpleMessage("True"))
	} else {
		this.sseServer.SendMessage("/api/event", sse.SimpleMessage("False"))
	}

	render.JSON(w, r, nil)
}
