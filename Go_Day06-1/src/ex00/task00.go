package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))

	// Set the background color to a bright blue
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{135, 206, 235, 255}}, image.ZP, draw.Src)

	// Draw Patrick's body
	patrickBodyColor := color.RGBA{255, 165, 0, 255} // orange
	patrickBodyRect := image.Rect(100, 100, 200, 200)
	draw.Draw(img, patrickBodyRect, &image.Uniform{patrickBodyColor}, image.ZP, draw.Src)

	// Draw Patrick's eyes
	patrickEyeColor := color.RGBA{0, 0, 0, 255} // black
	patrickEyeRect := image.Rect(130, 140, 150, 160)
	draw.Draw(img, patrickEyeRect, &image.Uniform{patrickEyeColor}, image.ZP, draw.Src)
	patrickEyeRect = image.Rect(160, 140, 180, 160)
	draw.Draw(img, patrickEyeRect, &image.Uniform{patrickEyeColor}, image.ZP, draw.Src)

	// Encode the image to PNG
	f, err := os.Create("image/png.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}
