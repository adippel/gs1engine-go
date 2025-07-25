# Example - GS1 Barcode Parser CLI

Example application for decoding **2D barcodes** from image files and reading its GS1 data:

- Supports decoding from **DataMatrix** and **Code128** barcodes
  using [makiuchi-d/gozxing](https://github.com/makiuchi-d/gozxing)
- Parses **GS1 syntax** using [adippel/gs1engine-go](https://github.com/adippel/gs1engine-go) syntax engine
- Works with image formats: PNG and JPEG

Scan performance of image files depends on [makiuchi-d/gozxing](https://github.com/makiuchi-d/gozxing) and its
configuration. Check out the library for more details.

## Usage

```bash
go run ./main.go -image path/to/image.png -type DATAMATRIX
```

### Options

| Flag     | Description                                    | Default      |
|----------|------------------------------------------------|--------------|
| `-image` | Path to the image file (png and jpg supported) | *(required)* |
| `-type`  | Barcode type: `DATAMATRIX` or `CODE128`        | `DATAMATRIX` |

## Example Output

Use `datamatrix-example.png` to test symbology `DataMatrix` together with `BarcodeMessageFormat`.

```bash
go run ./main.go -image ./datamatrix-example.png -type datamatrix
```

```text
2025/07/25 09:56:51 Successfully parsed GS message using syntax: BarcodeMessageFormat
2025/07/25 09:56:51 Found GS1 formatted message: (01)04012345123456
```

Use `code128-example.jpg` to test symbology `Code128` together with `BarcodeMessageScanData`.

```bash
go run ./main.go -image ./code128-example.jpg -type code128 
```

```text
2025/07/25 10:01:39 Parsed syntax type: BarcodeMessageScanData
2025/07/25 10:01:39 Detected symbology: ]C1
2025/07/25 10:01:39 Parsed GS1 message: (01)01234567890128(15)057072
```