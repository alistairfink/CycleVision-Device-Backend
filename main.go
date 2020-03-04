package main

import (
	"fmt"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/alistairfink/CycleVision-Device-Backend/Controllers"
	"github.com/alistairfink/CycleVision-Device-Backend/Utilities"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"os"
	"strings"
)

var gstreamPipeline string = "nvarguscamerasrc ! video/x-raw(memory:NVMM), width=1280, height=720, format=(string)NV12, framerate=21/1 ! nvvidconv flip-method=2 ! video/x-raw, width=1280, height=720, format=(string)BGRx ! videoconvert ! video/x-raw, format=(string)BGR ! appsink"

func main() {
	// Read Config
	var config *Utilities.Config
	if _, err := os.Stat("./Config.json"); err == nil {
		config = Utilities.GetConfig(".", "Config")
	} else {
		config = Utilities.GetConfig("/.", "Config")
	}

	// Router
	router, cleanup := routes()
	defer cleanup()
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf(" %-10s%-10s\n", method, strings.Replace(route, "/*", "", -1))
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

func routes() (*chi.Mux, func()) {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		CorsMiddleware,
	)

	// Controllers
	sseController := sse.NewServer(nil)
	signallingController := Controllers.NewSignallingController(sseController)
	videoStreamController, videoClose := Controllers.NewVideoStreamController(gstreamPipeline)
	go videoStreamController.Getframes()

	router.Route("/api", func(routes chi.Router) {
		routes.Mount("/event", sseController)
		routes.Mount("/signal", signallingController.Routes())
		routes.Mount("/video", videoStreamController.Routes())
	})

	return router, func() {
		sseController.Shutdown()
		videoClose()
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
