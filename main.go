package main

import (
	"fmt"
	imageprocessing "goroutines_pipeline/image_processing"
	"image"
	"os"
	"strings"
	"time"
)

type Job struct {
	InputPath string
	Image     image.Image
	OutPath   string
}

func loadImage(paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input path create a job and add it to
		// the out channel
		for _, p := range paths {
			job := Job{InputPath: p,
				OutPath: strings.Replace(p, "images/", "images/output/", 1)}
			job.Image = imageprocessing.ReadImage(p)
			out <- job
		}
		close(out)
	}()
	return out
}

func resize(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		// For each input job, create a new job after resize and add it to
		// the out channel
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Resize(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func convertToGrayscale(input <-chan Job) <-chan Job {
	out := make(chan Job)
	go func() {
		for job := range input { // Read from the channel
			job.Image = imageprocessing.Grayscale(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

func saveImage(input <-chan Job) <-chan bool {
	out := make(chan bool)
	go func() {
		for job := range input { // Read from the channel
			imageprocessing.WriteImage(job.OutPath, job.Image)
			out <- true
		}
		close(out)
	}()
	return out
}

func main() {
	var args = os.Args
	var usingGoroutines bool = false
	if len(args) >= 2 {
		if strings.TrimSpace(strings.ToLower(args[1])) == "goroutines" {
			usingGoroutines = true
		} else {
			fmt.Println("Unrecognised parameter!")
		}
	}

	imagePaths := []string{"images/image1.jpeg",
		"images/image2.jpeg",
		"images/image3.jpeg",
		"images/image4.jpeg",
	}

	runPipeline(imagePaths, usingGoroutines)
}

func runPipeline(imagePaths []string, useGoroutines bool) {
	if useGoroutines {
		fmt.Println("Running with goroutines:")
	} else {
		fmt.Println("Running without goroutines:")
	}
	start := time.Now()

	if useGoroutines {
		channel1 := loadImage(imagePaths)
		channel2 := resize(channel1)
		channel3 := convertToGrayscale(channel2)
		writeResults := saveImage(channel3)

		for success := range writeResults {
			if success {
				fmt.Println("Success!")
			} else {
				fmt.Println("Failed!")
			}
		}
	} else {
		for _, path := range imagePaths {
			job := Job{InputPath: path, OutPath: strings.Replace(path, "images/", "images/output/", 1)}
			img := imageprocessing.ReadImage(path)
			job.Image = img

			job.Image = imageprocessing.Resize(job.Image)
			job.Image = imageprocessing.Grayscale(job.Image)

			imageprocessing.WriteImage(job.OutPath, job.Image)
			fmt.Println("Success!")
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Pipeline completed in %s\n", elapsed)
}
