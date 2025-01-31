package main

import (
	"fmt"
	"log"

	"gocv.io/x/gocv"
)

func main() {
	inputFile := "input.MP4"
	outputFile := "output_clip.AVI" // Utilisez AVI au lieu de MP4

	// Ouvrir la vidéo
	video, err := gocv.VideoCaptureFile(inputFile)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la vidéo : %v", err)
	}
	defer video.Close()

	// Obtenir les propriétés de la vidéo
	frameRate := video.Get(gocv.VideoCaptureFPS)
	frameCount := int(video.Get(gocv.VideoCaptureFrameCount))
	width := int(video.Get(gocv.VideoCaptureFrameWidth))
	height := int(video.Get(gocv.VideoCaptureFrameHeight))

	fmt.Printf("Nombre total de frames dans la vidéo : %d\n", frameCount)
	fmt.Printf("Taux de frames (FPS) : %.2f\n", frameRate)
	fmt.Printf("Résolution de la vidéo : %dx%d\n", width, height)

	// Demander l'intervalle de frames
	fmt.Println("Entrez l'intervalle de frames à couper (par exemple, 10 40) :")
	var startFrame, endFrame int
	_, err = fmt.Scanf("%d %d", &startFrame, &endFrame)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture des entrées : %v", err)
	}

	// Ajuster les frames pour un index basé sur zéro
	startFrame--
	endFrame--

	// Valider l'intervalle de frames
	if startFrame < 0 || endFrame >= frameCount || startFrame >= endFrame {
		log.Fatalf("Intervalle de frames invalide. Veuillez entrer des valeurs correctes.")
	}

	// Créer un VideoWriter pour écrire la vidéo de sortie
	writer, err := gocv.VideoWriterFile(
		outputFile,
		"MJPG", // Codec MJPG pour AVI
		frameRate,
		width,
		height,
		true, // Couleur (true pour la vidéo couleur)
	)
	if err != nil {
		log.Fatalf("Erreur lors de la création du VideoWriter : %v", err)
	}
	defer writer.Close()

	// Lire et écrire les frames sélectionnées
	frame := gocv.NewMat()
	defer frame.Close()

	for i := 0; i < frameCount; i++ {
		if !video.Read(&frame) {
			log.Printf("Avertissement : Impossible de lire la frame %d\n", i)
			continue
		}

		// Si la frame est dans l'intervalle souhaité, l'écrire dans la vidéo de sortie
		if i >= startFrame && i <= endFrame {
			writer.Write(frame)
		}

		// Afficher la progression
		if i%10 == 0 {
			fmt.Printf("Traitement de la frame %d/%d\n", i, frameCount)
		}
	}

	fmt.Println("Découpage de la vidéo terminé avec succès.")
}
