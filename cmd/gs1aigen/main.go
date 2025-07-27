/*
Gs1aigen generates Go Code for GS1 AI Descriptions using the official [GS1 Syntax Dictionary].
For each AI description, it generates a corresponding [gs1.AIDescription].
In addition, a lookup table is generated mapping AI identifiers to their [gs1.AIDescription] (see the [gs1.AIRegistry]
for an example).

Usage:

	genairegistry [flags]

The flags are:

	-out string
	      Path to the output file (default "airegistry.go")
	-package string
	      Package name to use (default "gs1")
	-release string
	      Syntax Dictionary release to use (default "2025-01-30")

The release flag uses the Git Tags used in the [GS1 Syntax Dictionary] project.

[GS1 Syntax Dictionary]: https://github.com/gs1/gs1-syntax-dictionary
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/adippel/gs1engine-go"
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
