# AI Registry Generator CLI

The basis for working with GS1 data are the [Application Identifiers](https://ref.gs1.org/ai/). Those are standardized
together with an extensive description at [GS1 Syntax Dictionary](https://github.com/gs1/gs1-syntax-dictionary).

This CLI generates Go structures from it. This project's AI registry is also generated like this (
see [airegistry.go](../../airegistry.go))

The CLI can be used and configured as described next:

```bash
Usage of genairegistry:
  -out string
        Path to the output file (default "airegistry.go")
  -package string
        Package name to use (default "gs1engine")
  -release string
        Syntax Dictionary release to use (default "2025-01-30")
```

Run it as follows:

```bash
go run cmd/genairegistry/main.go -out ./myregistry.go
```