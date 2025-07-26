package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/adippel/gs1engine-go"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"github.com/makiuchi-d/gozxing/oned"
)

type barcodeType string

const (
	dataMatrix barcodeType = "DATAMATRIX"
	code128    barcodeType = "CODE128"
)

type options struct {
	ImagePath   string
	BarcodeType string
}

func Run(path string, barcodeType barcodeType) error {
	if path == "" {
		return errors.New("image path is empty")
	}

	// open and decode image file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening 2d image file: %w", err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding 2d image file: %w", err)
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return fmt.Errorf("error generating bitmap: %w", err)
	}

	var reader gozxing.Reader
	switch barcodeType {
	case dataMatrix:
		reader = datamatrix.NewDataMatrixReader()
	case code128:
		reader = oned.NewCode128Reader()
	default:
		return fmt.Errorf("unsupported barcode type: %s", barcodeType)
	}

	result, err := reader.Decode(
		bmp,
		map[gozxing.DecodeHintType]any{
			gozxing.DecodeHintType_TRY_HARDER:    true,
			gozxing.DecodeHintType_ASSUME_GS1:    true,
			gozxing.DecodeHintType_ALSO_INVERTED: true,
		},
	)
	if err != nil {
		log.Fatal("error decoding data matrix:", err)
	}

	msg, err := gs1.ParseMessage(result.GetText())
	if err != nil {
		log.Fatal("error parsing data message:", err)
	}

	log.Println("Parsed syntax type:", msg.SyntaxType)
	if msg.SyntaxType == gs1.BarcodeMessageScanData {
		log.Println("Detected symbology:", fmt.Sprintf("]%s%d", msg.Symbology.Type, msg.Symbology.Mode))
	}
	log.Println("Parsed GS1 message:", msg.AsElementString())

	return nil
}

func main() {
	var opts options
	flag.StringVar(&opts.ImagePath, "image", "", "path to image file")
	flag.StringVar(&opts.BarcodeType, "type", string(dataMatrix), "barcode type: oneOf(dataMatrix, code128)")
	flag.Parse()

	opts.BarcodeType = strings.ToUpper(opts.BarcodeType)

	err := Run(opts.ImagePath, barcodeType(opts.BarcodeType))
	if err != nil {
		log.Fatal(err)
	}
}
