package main

import (
	"fmt"
	"image/jpeg"
	"net/http"
	"os"

	"gocv.io/x/gocv"
)

func main() {
	vid, err := gocv.VideoCaptureFile("test.MP4")

	if err != nil {
		fmt.Printf("Error opening video capture: %v\n", err)
		return
	}
	defer vid.Close()

	//definir le gestionnaire HTTP pour la video
	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		//lirre une image de la video
		//img := image.NewRGBA(image.Rect(0, 0, 640, 480))
		timg := gocv.NewMat()
		ok := vid.Read(&timg)
		test := gocv.IMWrite("picture.jpeg", timg)
		if !test {
			fmt.Printf("Error of saving picture\n")
			return
		}
		imgff, err1 := os.Open("picture.jpeg")
		if err1 != nil {
			fmt.Printf("Error reading frame from video 1\n")
			return
		}
		img, err2 := jpeg.Decode(imgff)
		if err2 != nil {
			fmt.Printf("Error reading image from file \n")
			return
		}
		if !ok {
			fmt.Printf("Error reading frame from video 2\n")
			return
		}
		//convertir une image de type gocv en image

		//Encoder l'image et l'envoyer au client
		w.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(w, img, nil)
	})

	fmt.Println("Demarrage du serveur web sur http://localhost:8080/video")
	http.ListenAndServe(":8080", nil)
}
