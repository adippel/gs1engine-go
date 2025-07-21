package main

import (
	"flag"
	"fmt"
	"github.com/adippel/gs1engine-go"
	"log"
	"os"
)

type cliOpts struct {
	OutFile                 string
	SyntaxDictionaryRelease string
}

const defaultSyntaxDictionaryRelease = "2025-01-30"

func main() {
	var genOpts gs1.GenOpts
	var cliOpts cliOpts

	flag.StringVar(&genOpts.PackageName, "package", "gs1engine", "Package name to use")
	flag.StringVar(&cliOpts.OutFile, "out", "airegistry.go", "Path to the output file")
	flag.StringVar(&cliOpts.SyntaxDictionaryRelease, "release", defaultSyntaxDictionaryRelease, "Syntax Dictionary release to use")
	flag.Parse()

	data, err := gs1.DownloadSyntaxDictionary(cliOpts.SyntaxDictionaryRelease)
	if err != nil {
		panic(err)
	}

	ais, error := gs1.ParseSyntaxDictionary(&data)
	if error != nil {
		panic(error)
	}

	f, err := os.Create(cliOpts.OutFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = gs1.GenerateAIRegistry(f, ais, genOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Wrote %s\n", cliOpts.OutFile)
}
