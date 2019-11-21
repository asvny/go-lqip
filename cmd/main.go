package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-lqip/pkg/lqip"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	Version = "1.0.0"
)

var (
	inputFile           string
	version, outputJSON bool
)

type ImageData struct {
	Height             int               `json:"height"`
	Width              int               `json:"width"`
	PreviewSrc         string            `json:"previewSrc"`
	PreviewEnhancedSrc string            `json:"previewEnhancedSrc"`
	AspectRatio        float64           `json:"aspectRatio"`
	ColorPalette       map[string]string `json:"colorPalette"`
}

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, Version))
	}

	flag.StringVar(&inputFile, "i", "", "Filepath of the input file")
	flag.BoolVar(&outputJSON, "json", false, "Output JSON file")
	flag.BoolVar(&version, "version", false, "Show version and exit")
	flag.BoolVar(&version, "v", false, "Show version and exit")

	flag.Parse()
}

func main() {
	if version {
		fmt.Println(Version)
		return
	}

	// TODO : remove it
	if inputFile == "" {
		inputFile = "/Volumes/Personal/workspace/go-lqip/test-images/test.png"
	}

	lqipData, err := lqip.NewImage(inputFile)
	failOnErr(err)

	height, width := lqipData.Dimensions()
	aspectRatio := lqipData.AspectRatio()

	colorPaletteChan := make(chan map[string]string)
	go func() {
		colorPalette := make(map[string]string)
		_colorPalette, err := lqipData.ColorPalette()
		failOnErr(err)

		for name, color := range _colorPalette {
			colorPalette[name] = fmt.Sprintf("#%s", strconv.Itoa(int(color)))
		}

		colorPaletteChan <- colorPalette
	}()

	base64Chan := make(chan string)
	go func() {
		base64, err := lqipData.PreviewSrc()
		failOnErr(err)

		base64Chan <- base64
	}()

	base64EnhancedChan := make(chan string)
	go func() {
		base64Enhanced, err := lqipData.PreviewEnhancedSrc()
		failOnErr(err)

		base64EnhancedChan <- base64Enhanced
	}()

	imageData := ImageData{
		Height:             height,
		Width:              width,
		AspectRatio:        aspectRatio,
		PreviewSrc:         <-base64Chan,
		PreviewEnhancedSrc: <-base64EnhancedChan,
		ColorPalette:       <-colorPaletteChan,
	}

	if outputJSON {
		paths := strings.Split(inputFile, ".")
		err := writeAsJSONFile(paths[0]+".json", imageData)
		failOnErr(err)

		return
	}

	printTabularView(imageData)
}

func writeAsJSONFile(outputFilePath string, imgData ImageData) error {
	json, _ := json.MarshalIndent(imgData, "", "  ")

	err := ioutil.WriteFile(outputFilePath, json, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func printTabularView(img ImageData) {
	colors := ""
	for name, color := range img.ColorPalette {
		colors = colors + fmt.Sprintf("%s - %s\n", name, color)
	}

	data := [][]string{
		[]string{"Height", fmt.Sprintf("%d", img.Height)},
		[]string{"Width", fmt.Sprintf("%d", img.Width)},
		[]string{"Aspect ratio", fmt.Sprintf("%f", img.AspectRatio)},
		[]string{"Color Pallette", colors},
		[]string{"Preview src", hardWrap(img.PreviewSrc, 120)},
		[]string{"Preview enhanced src", hardWrap(img.PreviewEnhancedSrc, 120)},
	}

	table := tablewriter.NewWriter(os.Stdout)

	fmt.Println("File ::: ", inputFile)

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}

/// Util ///

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func hardWrap(text string, colBreak int) string {
	if colBreak < 1 {
		return text
	}

	text = strings.TrimSpace(text)
	wrapped := ""

	var i int
	for i = 0; len(text[i:]) > colBreak; i += colBreak {
		wrapped += text[i:i+colBreak] + "\n"
	}

	wrapped += text[i:]

	return wrapped
}

const usage = `
lqip %s
A cli tool to generate Low Quality Image Placeholder(LQIP) and 
it outputs aspect-ratio, image size and encoded base64 lqip images.

lqip [options...] <interface>
Options:
  -version, -v		print version and exit
  -i			input path of the image file
  -o			output path of the json file (if empty it uses input file name)                 
`
