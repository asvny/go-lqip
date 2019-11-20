package main

import (
	"flag"
	"fmt"
	"go-lqip/pkg/lqip"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	inputFile      = flag.String("i", "/Volumes/Personal/workspace/go-lqip/test-images/test.png", "Filepath of the input file")
	outputFileName = flag.String("o", "output.gif", "Name of the output file")
	quality        = flag.Int("q", 5, "Quality of the placeholder image")
	imageType      = flag.String("t", "jpeg", "Type of the base64 output image")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprint(usage, 1.5, runtime.NumCPU()))
	}

	flag.Parse()

	imageFile, err := readImageFile(*inputFile)
	failOnErr(err)

	lqipImage := lqip.NewImage(imageFile)
	height, width := lqipImage.Size()
	aspectRatio := lqipImage.AspectRatio()
	colorPalette := lqipImage.ColorPalette()

	colors := ""
	for name, color := range colorPalette {
		colors = colors + fmt.Sprintf("%s - #%d\n", name, color)
	}

	base64, err := lqipImage.PreviewSrc()
	failOnErr(err)

	base64Enhanced, err := lqipImage.PreviewEnhancedSrc()
	failOnErr(err)

	data := [][]string{
		[]string{"Height", fmt.Sprintf("%d", height)},
		[]string{"Width", fmt.Sprintf("%d", width)},
		[]string{"Aspect ratio", fmt.Sprintf("%f", aspectRatio)},
		[]string{"Color Pallette", colors},
		[]string{"Preview src", hardWrap(base64, 120)},
		[]string{"Preview enhanced src", hardWrap(base64Enhanced, 120)},
	}

	table := tablewriter.NewWriter(os.Stdout)

	fmt.Println("File ::: ", *inputFile)

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}

func readImageFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
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
DEMO Content
`
