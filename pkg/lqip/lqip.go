package lqip

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/generaltso/vibrant"
)

// Web safe color in HEX format
type Color vibrant.Color

type ImageOps interface {
	// Returns the height/width of the image
	Size() (int, int)
	// Returns the aspect ratio of the image
	AspectRatio() float64
	// Returns a low quality placeholder base64 string
	PreviewSrc() string
	// Returns the color palette of the image
	ColorPalette() map[string]Color
}

func NewImage(imageFile *os.File) *Image {
	return &Image{
		file:       imageFile,
		hasDecoded: false,
	}
}

// Image has basic file properties of the image
type Image struct {
	file        *os.File
	fileBuffer  *bytes.Buffer
	imageFile   image.Image
	hasDecoded  bool
	imageConfig image.Rectangle
	format      string
}

func (i *Image) Size() (int, int) {
	if !(i.hasDecoded) {
		i.storeImageConfig()
	}

	return i.imageConfig.Dy(), i.imageConfig.Dx()
}

func (i *Image) AspectRatio() float64 {
	if !(i.hasDecoded) {
		i.storeImageConfig()
	}

	return toFixed(float64(i.imageConfig.Dx()) / float64(i.imageConfig.Dy()))
}

// Extracts colors from the image
func (i *Image) ColorPalette() map[string]Color {
	palette, err := vibrant.NewPalette(i.imageFile, 32)
	if err != nil {
		log.Fatal(err)
	}

	colorMap := make(map[string]Color)

	for name, swatch := range palette.ExtractAwesome() {
		colorMap[name] = Color(swatch.Color)
	}

	return colorMap
}

func (i *Image) PreviewSrc() (string, error) {
	base64, err := i.resizeAndBase64(3, 3)
	return base64, err
}

func (i *Image) PreviewEnhancedSrc() (string, error) {
	base64, err := i.resizeAndBase64(12, 12)
	return base64, err
}

func (i *Image) resizeAndBase64(width, height int) (string, error) {
	resizedImage := imaging.Resize(i.imageFile, width, height, imaging.Lanczos)

	base64String, err := imageToBase64(resizedImage, i.format)
	if err != nil {
		return "", err
	}

	return base64String, nil
}

func (i *Image) storeImageConfig() {
	image, format, err := image.Decode(i.file)

	if err != nil {
		log.Fatal(err)
	}

	i.imageConfig = image.Bounds()
	i.imageFile = image
	i.format = format
	i.hasDecoded = true
}

func toFixed(num float64) float64 {
	return float64(int(num*100)) / 100
}

func imageToBase64(img image.Image, format string) (string, error) {
	var buff bytes.Buffer
	var err error

	switch format {
	case "png":
		png.Encode(&buff, img)
	case "gif":
		gif.Encode(&buff, img, nil)
	case "jpeg":
		jpeg.Encode(&buff, img, nil)
	default:
		err = errors.New("Unsupported image format")
	}

	if err != nil {
		return "", err
	}

	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())
	htmlImage := fmt.Sprintf("data:image/%s;base64,%s", format, encodedString)

	return htmlImage, nil
}
