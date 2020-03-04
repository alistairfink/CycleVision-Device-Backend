package Controllers

import (
	"fmt"
	"github.com/go-chi/chi"
	"gocv.io/x/gocv"
	"image"
	"net/http"
	"sync"
	"time"
)

type VideoStreamController struct {
	camera  *gocv.VideoCapture
	frame   []byte
	mutex   *sync.Mutex
	frameId int
}

func NewVideoStreamController(gstreamPipeline string) (*VideoStreamController, func()) {
	controller := VideoStreamController{
		camera: gocv.OpenVideoCapture(gstreamPipeline),
		frame:  []byte{},
		mutex:  &sync.Mutex{},
	}

	return &controller, func() {
		controller.camera.Close()
	}
}

func (this *VideoStreamController) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", this.Get)
	return router
}

func (this *VideoStreamController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	data := ""
	for {
		this.mutex.Lock()
		data = "--frame\r\n  Content-Type: image/jpeg\r\n\r\n" + string(this.frame) + "\r\n\r\n"
		this.mutex.Unlock()
		time.Sleep(33 * time.Millisecond)
		w.Write([]byte(data))
	}
}

func (this *VideoStreamController) Getframes() {
	img := gocv.NewMat()
	defer img.Close()
	for {
		if ok := this.camera.Read(&img); !ok {
			fmt.Printf("Device closed\n")
			return
		}
		if img.Empty() {
			continue
		}
		this.frameId++
		gocv.Resize(img, &img, image.Point{}, float64(0.5), float64(0.5), 0)
		frame, _ = gocv.IMEncode(".jpg", img)

	}
}
