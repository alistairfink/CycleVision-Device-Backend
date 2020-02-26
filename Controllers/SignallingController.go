package Controllers

import (
	"encoding/json"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"io/ioutil"
	"log"
	"net/http"
)

type Signal struct {
	Signal bool
}

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
	router.Post("/", this.Post)
	return router
}

func (this *SignallingController) Post(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Print(err)
		http.Error(w, "Error Processing Request", http.StatusBadRequest)
		return
	}

	var signalModel Signal
	err = json.Unmarshal(body, &signalModel)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error Processing Model", http.StatusBadRequest)
		return
	}

	if signalModel.Signal {
		this.sseServer.SendMessage("/api/event", sse.SimpleMessage("True"))
	} else {
		this.sseServer.SendMessage("/api/event", sse.SimpleMessage("False"))
	}

	render.JSON(w, r, nil)
}
