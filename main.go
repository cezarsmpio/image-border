package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func convertHexToRgba(hex string) color.RGBA {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return color.RGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}

func main() {
	// Get the folder path from command-line arguments
	var folderPath string
	var padding int
	var borderColor string

	flag.StringVar(&folderPath, "folder", "", "The folder path containing the images")
	flag.IntVar(&padding, "padding", 200, "The padding size for the border")
	flag.StringVar(&borderColor, "color", "#ffffff", "The color for the border")

	flag.Parse()

	// Create the "results" folder inside the specified path
	resultsFolderPath := filepath.Join(folderPath, "results")
	err := os.MkdirAll(resultsFolderPath, 0755)
	if err != nil {
		fmt.Println("Failed to create results folder:", err)
		return
	}

	// Read all image files from the specified folder
	imageFiles, err := filepath.Glob(filepath.Join(folderPath, "*.jpg"))
	if err != nil {
		fmt.Println("Failed to read image files:", err)
		return
	}

	// Process each image file
	wg := sync.WaitGroup{}
	for _, imageFile := range imageFiles {
		wg.Add(1)
		fmt.Println("Processing image:", imageFile)

		go func() {
			defer wg.Done()

			// Open the image file
			file, err := os.Open(imageFile)
			if err != nil {
				fmt.Println("Failed to open image file:", err)
				return
			}
			defer file.Close()

			// Decode the image
			img, _, err := image.Decode(file)
			if err != nil {
				fmt.Println("Failed to decode image:", err)
				return
			}

			// Create a new image with the white border
			borderColor := convertHexToRgba(borderColor)
			bounds := img.Bounds()
			borderWidth := bounds.Dx() + 2*padding
			borderHeight := bounds.Dy() + 2*padding
			borderImg := image.NewRGBA(image.Rect(0, 0, borderWidth, borderHeight))
			draw.Draw(borderImg, borderImg.Bounds(), &image.Uniform{borderColor}, image.Point{}, draw.Src)
			draw.Draw(borderImg, image.Rect(padding, padding, padding+bounds.Dx(), padding+bounds.Dy()), img, bounds.Min, draw.Src)

			// Create the output file path
			outputFilePath := filepath.Join(resultsFolderPath, filepath.Base(imageFile))

			// Save the modified image to the output file
			outputFile, err := os.Create(outputFilePath)
			if err != nil {
				fmt.Println("Failed to create output file:", err)
				return
			}
			defer outputFile.Close()

			err = jpeg.Encode(outputFile, borderImg, nil)
			if err != nil {
				fmt.Println("Failed to encode output image:", err)
				return
			}

			fmt.Println("Modified image saved:", outputFilePath)
		}()
	}

	wg.Wait()
}
