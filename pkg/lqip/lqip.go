package lqip

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

type ImageOps interface {
	// Returns the height/width of the image
	Size() (int, int)
	// Returns the aspect ratio of the image
	AspectRatio() float64
	// Returns a low quality placeholder base64 string
	PreviewSrc() string
	// Returns the image URL
	Src() string
	// Returns the color palette of the image
	ColorPalette() []string
}

// Image has basic file properties of the image
type Image struct {
	file       *os.File
	fileBuffer *bytes.Buffer
	imageFile  image.Image
	hasDecoded bool

	imageConfig image.Config
}

func (i *Image) Size() (int, int) {
	if !(i.hasDecoded) {
		i.storeImageConfig()
	}

	return i.imageConfig.Height, i.imageConfig.Width
}

func (i *Image) AspectRatio() float64 {
	if !(i.hasDecoded) {
		i.storeImageConfig()
	}

	return toFixed(float64(i.imageConfig.Height) / float64(i.imageConfig.Width))
}

func (i *Image) Src() string {
	return ""
}

func (i *Image) ColorPalette() []string {
	return []string{"#000"}
}

func (i *Image) storeImageConfig() {
	image, _, err := image.DecodeConfig(i.file)
	if err != nil {
		log.Fatal(err)
	}

	i.imageConfig = image
	i.hasDecoded = true
}

func NewImage(imageFile *os.File) *Image {
	return &Image{
		file:       imageFile,
		hasDecoded: false,
	}
}

func toFixed(num float64) float64 {
	return float64(int(num*100)) / 100
}
