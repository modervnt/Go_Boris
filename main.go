package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"time"

	"gocv.io/x/gocv"
)

var (
	err      error
	webcam   *gocv.VideoCapture
	frame_id int
)

func main() {
	host := "localhost:3000"

	// open webcam
	if len(os.Args) < 2 {
		fmt.Println(">> device /dev/video0 (default)")
		webcam, err = gocv.VideoCaptureDevice(0)
	} else {
		fmt.Println(">> file/url :: " + os.Args[1])
		webcam, err = gocv.VideoCaptureFile(os.Args[1])
	}

	if err != nil {
		fmt.Printf("Error opening capture device: \n")
		return
	}
	defer webcam.Close()

	// Servir les fichiers statiques (HTML, CSS, JS)
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Route pour la vidéo
	http.HandleFunc("/video", videoHandler)

	fmt.Printf("Server started at http://%s\n", host)
	log.Fatal(http.ListenAndServe(host, nil))
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed\n")
			return
		}
		if img.Empty() {
			continue
		}

		frame_id++
		gocv.Resize(img, &img, image.Point{}, float64(0.5), float64(0.5), 0)
		frame, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			fmt.Printf("Error encoding frame: %v\n", err)
			continue
		}

		fmt.Println("Frame ID: ", frame_id)

		data := "--frame\r\nContent-Type: image/jpeg\r\n\r\n" + string(frame.GetBytes()) + "\r\n"
		_, err = w.Write([]byte(data))
		if err != nil {
			return
		}

		time.Sleep(33 * time.Millisecond)
	}
}
