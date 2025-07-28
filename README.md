# Go GS1 Barcode Syntax Engine

[![Go Reference](https://pkg.go.dev/badge/github.com/adippel/gs1engine-go.svg)](https://pkg.go.dev/github.com/adippel/gs1engine-go)

A pure Go implementation for working with GS1 Data in various syntax variations as well as with
[GS1 Application Identifiers (AI)](https://ref.gs1.org/ai/).

**Highlights:**

* ✅ Pure Go implementation, zero external dependencies
* ✅ Parser support for the following syntax types:
	* GS1 element string syntax (e.g. `(01)09526064055028(17)250521(10)ABC123(21)456DEF`)
	* Barcode message format (e.g. `^01095260640550281725052110ABC123^21456DEF`)
	* Barcode message scan data (e.g. `]d201095260640550281725052110ABC123{GS}21456DEF`)
* ✅ AI Registry with description of all 536 AIs (as of release `2025-01-30`)
* ✅ [Go Code generator CLI](./cmd/gs1aigen/README.md) to generate AI description based on the official
  [GS1 Syntax Dictionary](https://github.com/gs1/gs1-syntax-dictionary)
* ✅ Usable examples in `examples/`

**⛔️ NO STABLE API yet.** Not functional complete, see the following roadmap:

* Implement GS1 data parsing support for:
	* Digital-Link URI Syntax (e.g. `https://example.com/01/09526064055028`)
* Implement high-level interface to easily work with GS1 description (Flags, Specification, Attributes)
* Implement validation support by implementing linters to validate that an AI conforms
	* to its Specification (e.g `yymmdd` or `N6`)
	* and to its attributes (e.g `req` and `ex` to define valid and invalid pairings).
* Add GitHub workflows:
	* Unit Testing with code coverage
	* Linting
* Add badges for build status and code coverage

## Usage

### Installing

Add it to your project:

```bash
go get github.com/adippel/gs1engine-go
```

### Examples

Several examples exist in `examples/` showcasing the parser:

* [read2dcode - Example to parse GS1 messages from 2D barcodes](./examples/read2dcode/README.md)
* [cliparser - Example to parse GS1 message from CLI input](./examples/cliparser/README.md)

### Parsing

Several GS1 syntax formats are supported via the following functions:

- [ParseMessage](https://pkg.go.dev/github.com/adippel/gs1engine-go#ParseMessage): Automatically detects syntax type
  based on input and then uses the correct parsers.
- [ParseBarcodeMessage](https://pkg.go.dev/github.com/adippel/gs1engine-go#ParseBarcodeMessage): Parses barcode message
  format and barcode scan data (e.g. `]d2...`, `^...`)
- [ParseElementString](https://pkg.go.dev/github.com/adippel/gs1engine-go#ParseElementString): Parses element string
  syntax (e.g `(01)...(17)...`)

All parser support the visual FNC1 substitutes `^` and `{GS}`.

🛑 Plain syntax (non-AI form) is not supported.

To start parsing, use the following:

```go
gs1Messages := []string{
	"]d201095260640550281725052110ABC123{GS}21456DEF",
	"^01095260640550281725052110ABC123^21456DEF",
	"(01)09526064055028(17)250521(10)ABC123(21)456DEF",
}

for _, gs1Msg := range gs1Messages {
	gs1Data, _ := gs1.ParseMessage(gs1Msg)
	fmt.Println("Detected GS1 Syntax Type:", gs1Data.SyntaxType)
	fmt.Println("Message:", gs1Data.AsElementString())
}
```

## References

The GS1 has good reference material to understand their system and approaches:

* To see the whole collection of available documents, see the [GS1 reference](https://ref.gs1.org)
* For technical details,
  see [GS1 Barcode Syntax Resource User Guide](https://ref.gs1.org/tools/gs1-barcode-syntax-resource/user-guide/)
* For a broader picture,
  see [GS1 System Architecture](https://www.gs1.org/standards/gs1-system-architecture-document/current-standard)
* [GS1-maintained syntax engine in C](https://github.com/gs1/gs1-syntax-engine)