# Example - GS1 CLI Parser

This is a simple command-line tool that demonstrates how to parse GS1 data using the
[adippel/gs1engine-go](https://github.com/adippel/gs1engine-go) library.

It uses `ParseDataMessage` to automatically detect and parse the GS1 syn@tax type of the input data.

## Usage

Build or run the program and pass the GS1 message via the `-message` flag:

```bash
go run main.go -message "]d201095260640550281725052110ABC123{GS}21456DEF"
```

**Example Inputs**

| Syntax Type            | Example Input                                      |
|------------------------|----------------------------------------------------|
| Barcode Scan Data      | `]d201095260640550281725052110ABC123{GS}21456DEF`  |
| Barcode Message Format | `^01095260640550281725052110ABC123^21456DEF`       |
| Element String Syntax  | `(01)09526064055028(17)250521(10)ABC123(21)456DEF` |