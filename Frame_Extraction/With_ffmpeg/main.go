package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	inputFile := "input.MP4"
	outputFile := "output_clip.MP4"

	cmdGet := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-count_frames", "-show_entries", "stream=nb_read_frames,avg_frame_rate", "-of", "csv=p=0", inputFile)

	var out bytes.Buffer
	cmdGet.Stdout = &out
	err := cmdGet.Run()
	if err != nil {
		log.Fatalf("Error of ffprobe execution : %v", err)
	}
	results := strings.Split(strings.TrimSpace(out.String()), ",")
	if len(results) != 2 {
		log.Fatalf("Error of outing format : %v", results)
	}
	frameCount, err := strconv.Atoi(results[1])
	if err != nil {
		log.Fatalf("Error of conversion : %v", err)
	}
	fmt.Printf("This video have %d frames \n", frameCount)

	// Conversion du taux de frames (FPS)
	frameRateParts := strings.Split(results[0], "/")
	if len(frameRateParts) != 2 {
		log.Fatalf("Unexpected frame rate format : %v", results[1])
	}
	numerator, err := strconv.Atoi(frameRateParts[0])
	if err != nil {
		log.Fatalf("Error of conversion : %v", err)
	}
	denominator, err := strconv.Atoi(frameRateParts[1])
	if err != nil {
		log.Fatalf("Error of conversion : %v", err)
	}
	frameRate := float64(numerator) / float64(denominator)

	/*frameRate, err := strconv.ParseFloat(results[1], 64)
	if err != nil {
		log.Fatal("Error of conversion : %v", err)
	}*/

	fmt.Println("Enter the interval of frames to be cut:")
	var startFrame, endFrame int
	_, err = fmt.Scanf("%d %d", &startFrame, &endFrame)
	if err != nil {
		log.Fatalf("Error of reading data : %v", err)
	}
	startFrame--
	endFrame--

	startTime := float64(startFrame) / frameRate
	endTime := float64(endFrame) / frameRate

	cmdCut := exec.Command("ffmpeg", "-i", inputFile, "-ss", fmt.Sprintf("%.2f", startTime), "-to", fmt.Sprintf("%.2f", endTime), "-c", "copy", outputFile)
	output, err := cmdCut.CombinedOutput()
	if err != nil {
		log.Fatalf("Error of ffmpeg cutting command : %v\n Out : %s", err, string(output))
	}

	fmt.Println("End!")
}
