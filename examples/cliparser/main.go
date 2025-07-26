package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/adippel/gs1engine-go"
)

func main() {
	msg := flag.String("message", "", "GS1 data to parse")
	flag.Parse()

	if *msg == "" {
		log.Fatal("You must specify a message to parse")
	}

	gs1Data, err := gs1.ParseMessage(*msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Parsed syntax type:", gs1Data.SyntaxType)
	if gs1Data.SyntaxType == gs1.BarcodeMessageScanData {
		log.Println("Detected symbology:", fmt.Sprintf("]%s%d", gs1Data.Symbology.Type, gs1Data.Symbology.Mode))
	}
	log.Println("Parsed GS1 message:", gs1Data.AsElementString())
}
