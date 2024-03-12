package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestConvertHexToRGBA(t *testing.T) {
	hex := "#ff0000"
	expected := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	result := convertHexToRgba(hex)

	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestMain(t *testing.T) {
	padding := 200
	borderColor := "#ffffff"

	// Create a temporary folder for test results
	tempFolderPath, err := os.MkdirTemp("", "test-folder")
	resultsFolderPath := filepath.Join(tempFolderPath, "results")
	if err != nil {
		t.Fatalf("Failed to create results folder: %v", err)
	}
	defer os.RemoveAll(tempFolderPath)

	// Create a temporary image file for testing
	tempImageFile := filepath.Join(tempFolderPath, "test.jpg")
	err = createTempImageFile(tempImageFile)
	if err != nil {
		t.Fatalf("Failed to create temporary image file: %v", err)
	}
	defer os.Remove(tempImageFile)

	// Run the main function
	os.Args = []string{"cmd", "-folder", tempFolderPath, "-padding", strconv.Itoa(padding), "-color", borderColor}
	fmt.Println(tempFolderPath, resultsFolderPath, tempImageFile)
	main()

	// Check if the modified image file was created
	outputFilePath := filepath.Join(resultsFolderPath, "test.jpg")
	_, err = os.Stat(outputFilePath)
	if os.IsNotExist(err) {
		t.Errorf("Modified image file was not created: %v", err)
	}
}

func createTempImageFile(filePath string) error {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create image file:", err)
		return err
	}
	defer file.Close()
	err = jpeg.Encode(file, img, nil)
	if err != nil {
		fmt.Println("Failed to encode image:", err)
		return err
	}
	return nil
}
