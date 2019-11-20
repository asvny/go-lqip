package main

import (
	"flag"
	"fmt"
	"go-lqip/pkg/lqip"
	"log"
	"os"
	"runtime"
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

	fmt.Println(lqipImage.Size())
	fmt.Println(lqipImage.AspectRatio())
	fmt.Println(lqipImage.ColorPalette())

	base64, err := lqipImage.PreviewSrc()
	failOnErr(err)

	fmt.Println(base64)

	base64Enhanced, err := lqipImage.PreviewEnhancedSrc()
	failOnErr(err)

	fmt.Println()
	fmt.Println(base64Enhanced)
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

const usage = `
DEMO Content
`
