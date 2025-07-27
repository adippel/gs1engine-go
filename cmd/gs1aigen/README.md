# GS1 AI Generator

This CLI generates Go Code for GS1 AI Descriptions using the
official [GS1 Syntax Dictionary](https://github.com/gs1/gs1-syntax-dictionary/blob/main/gs1-syntax-dictionary.txt).

The basis for working with GS1 data are the [Application Identifiers](https://ref.gs1.org/ai/). Those are standardized
together with an extensive description at [GS1 Syntax Dictionary](https://github.com/gs1/gs1-syntax-dictionary).

Usage:

    genairegistry [flags]

The flags are:

    -out string
          Path to the output file (default "airegistry.go")
    -package string
          Package name to use (default "gs1")
    -release string
          Syntax Dictionary release to use (default "2025-01-30")

The release flag uses the [Git Tags](https://github.com/gs1/gs1-syntax-dictionary/tags) used in the GS1 Syntax
Dictionary project.