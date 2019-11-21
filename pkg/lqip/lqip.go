package lqip

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/generaltso/vibrant"
	E "github.com/pkg/errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
)

// Meaningful error messages
var (
	ErrImageConfig  = "Cannot obtain image config from the file"
	ErrColorPalette = "Cannot obtain color palette from the image"
	ErrImageBase64  = "Cannot convert resized image to base64 string"
)

// Color - Web basded color format
type Color vibrant.Color

// ImageOps - List of public image functions form Image struct
type ImageOps interface {
	// Returns the height/width of the image
	Dimensions() (int, int)
	// Returns the aspect ratio of the image
	AspectRatio() float64
	// Returns a low quality placeholder base64 string
	PreviewSrc() (string, error)
	// Returns a low quality placeholder base64 string
	PreviewEnhancedSrc() (string, error)
	// Returns the color palette of the image
	ColorPalette() (map[string]Color, error)
}

// Image has basic file properties of the image
type Image struct {
	imagePath  string
	imageFile  image.Image
	hasDecoded bool
	format     string
}

// NewImage implements the ImageOps interface
func NewImage(imagePath string) (ImageOps, error) {
	image := &Image{
		imagePath:  imagePath,
		hasDecoded: false,
	}

	err := image.storeImageConfig()
	if err != nil {
		return nil, E.Wrap(err, ErrImageConfig)
	}

	return image, nil
}

// Dimensions returns the height/width of the image
func (i *Image) Dimensions() (int, int) {
	bounds := i.imageFile.Bounds()
	return bounds.Dy(), bounds.Dx()
}

// AspectRatio returns the aspect ratio of the image
func (i *Image) AspectRatio() float64 {
	height, width := i.Dimensions()
	return toFixed(float64(height) / float64(width))
}

// ColorPalette returns the extracted colors from the image
func (i *Image) ColorPalette() (map[string]Color, error) {
	palette, err := vibrant.NewPalette(i.imageFile, 32)
	if err != nil {
		return nil, E.Wrap(err, ErrColorPalette)
	}

	colorMap := make(map[string]Color)
	for name, swatch := range palette.ExtractAwesome() {
		colorMap[name] = Color(swatch.Color)
	}

	return colorMap, nil
}

// PreviewSrc returns base64-encoded image which can be used as placeholder for small dimension preview
func (i *Image) PreviewSrc() (string, error) {
	base64, err := i.resizeAndBase64(3, 3)
	return base64, E.Wrap(err, ErrImageBase64)
}

// PreviewEnhancedSrc returns base64-encoded image which can be used as placeholder for large dimension preview
func (i *Image) PreviewEnhancedSrc() (string, error) {
	base64, err := i.resizeAndBase64(12, 12)
	return base64, E.Wrap(err, ErrImageBase64)
}

func (i *Image) storeImageConfig() error {
	imageFile, format, err := readImageFile(i.imagePath)
	if err != nil {
		return err
	}

	i.format = format
	i.imageFile = imageFile

	return nil
}

func (i *Image) resizeAndBase64(width, height int) (string, error) {
	resizedImage := imaging.Resize(i.imageFile, width, height, imaging.Lanczos)

	base64String, err := imageToBase64(resizedImage, i.format)
	if err != nil {
		return "", err
	}

	return base64String, nil
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

func readImageFile(filePath string) (image.Image, string, error) {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return nil, "", err
	}

	imageFile, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	return imageFile, format, nil
}
