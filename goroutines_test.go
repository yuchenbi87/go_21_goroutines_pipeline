package main

import (
	imageprocessing "goroutines_pipeline/image_processing"
	"testing"
)

func TestLoadImage(t *testing.T) {
	paths := []string{"images/image1.jpeg", "images/image2.jpeg"}
	ch := loadImage(paths)

	count := 0
	for job := range ch {
		if job.Image == nil {
			t.Errorf("Failed to load image: %s", job.InputPath)
		}
		count++
	}

	if count != len(paths) {
		t.Errorf("Expected %d images, but got %d", len(paths), count)
	}
}

func TestResize(t *testing.T) {
	inputJob := Job{InputPath: "images/image1.jpeg"}
	inputJob.Image = imageprocessing.ReadImage(inputJob.InputPath)
	input := make(chan Job, 1)
	input <- inputJob
	close(input)

	output := resize(input)
	job := <-output

	if job.Image == nil {
		t.Errorf("Failed to resize image: %s", inputJob.InputPath)
	}
}

func TestConvertToGrayscale(t *testing.T) {
	inputJob := Job{InputPath: "images/image1.jpeg"}
	inputJob.Image = imageprocessing.ReadImage(inputJob.InputPath)
	input := make(chan Job, 1)
	input <- inputJob
	close(input)

	output := convertToGrayscale(input)
	job := <-output

	if job.Image == nil {
		t.Errorf("Failed to convert image to grayscale: %s", inputJob.InputPath)
	}
}

func TestSaveImage(t *testing.T) {
	inputJob := Job{InputPath: "images/image1.jpeg", OutPath: "images/output/image1_output.jpeg"}
	inputJob.Image = imageprocessing.ReadImage(inputJob.InputPath)
	input := make(chan Job, 1)
	input <- inputJob
	close(input)

	output := saveImage(input)
	success := <-output

	if !success {
		t.Errorf("Failed to save image: %s", inputJob.OutPath)
	}
}
