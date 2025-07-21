# Go GS1 Barcode Syntax Engine

A pure Go implementation for working with GS1 Data and Syntax inspired by the official GS1 C-reference implementation
[GS1 Barcode Syntax Engine](https://github.com/gs1/gs1-syntax-engine). This implementation is not feature-complete and
mature enough (yet). Thus it cannot provide you with the same guarantees. Though, it can be a great help and the basis
for working with GS1 data and [GS1 Application Identifiers (AI)](https://ref.gs1.org/ai/).

Start using the project by adding it to your project:

```bash
go get github.com/adippel/gs1engine-go
```

## Status Quo

**Features**

- Pure Go implementation, no CGo
- Generate Go data types for all AIs using the [GS1 Syntax Dictionary](https://github.com/gs1/gs1-syntax-dictionary)

**Roadmap**

* Implement GS1 data decoder supporting the following formats:
  * GS1 element string syntax (e.g. `(01)09526064055028(17)250521(10)ABC123(21)456DEF`)
  * Barcode message format (e.g. `^01095260640550281725052110ABC123^21456DEF`)
  * Barcode message scan data (e.g. `]d201095260640550281725052110ABC123{GS}21456DEF`)
* Implement high-level interface to easily work with GS1 description (Flags, Specification, Attributes)
* Implement linter-support to validate GS-AI types based on the flags

## Using the AI Code generator

The basis for working with GS1 data are the [Application Identifiers](https://ref.gs1.org/ai/). Those are standardized
together with an extensive description at [GS1 Syntax Dictionary](https://github.com/gs1/gs1-syntax-dictionary). 

To generate Go structures from it, use program `cmd/genairegistry` as follows:

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