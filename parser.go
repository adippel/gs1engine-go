package gs1

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	fnc1          = '\x1d' // fnc1 character used to terminate variable length AIs
	symbologyFlag = ']'    // symbologyFlag is the character used to indicate start of symbology information
)

// fnc1Visuals are visual representations commonly used to substitute the FNC1 char.
var fnc1Visuals = []string{"^", "{GS}"}

// ParseMessage detects the type of encoding used in msg and decodes the [Message] by dispatching to the more
// specialized parsers [ParseBarcodeMessage] and [ParseElementString].
// Plain syntax data is not parsed, as it is not AI-based.
func ParseMessage(msg string) (d Message, _ error) {
	if len(msg) == 0 {
		return d, errors.New("message is empty")
	}
	firstChar := msg[0]

	switch firstChar {
	case symbologyFlag, fnc1:
		return ParseBarcodeMessage(msg)
	case '(':
		return ParseElementString(msg)
	}
	for _, visualFNC1 := range fnc1Visuals {
		if strings.Contains(msg, visualFNC1) {
			return ParseBarcodeMessage(msg)
		}
	}

	return d, errors.New("unsupported data syntax")
}

// ParseBarcodeMessage supports `Barcode message format` and `Barcode message scan data`. It supports group
// separation with FNC1 as well as its literal variants '^' and '{GS}'. Examples for messages are
// “^01095260640550281725052110ABC123^21456DEF` and `]d201095260640550281725052110ABC123{GS}21456DEF`.
func ParseBarcodeMessage(msg string) (d Message, _ error) {
	// Clean the input string
	for _, visualFNC1 := range fnc1Visuals {
		msg = strings.ReplaceAll(msg, visualFNC1, string(fnc1))
	}
	if len(msg) > 0 && msg[0] == fnc1 {
		msg = msg[1:] // Remove leading FNC1
	}

	// Read and remove symbology identifier if present (e.g., ]C1)
	if strings.HasPrefix(msg, string(symbologyFlag)) && len(msg) > 3 {
		symbologyMode, err := strconv.Atoi(msg[2:3])
		if err != nil {
			return d, fmt.Errorf("error parsing symbology mode: %w", err)
		}
		d.Symbology.Type = SymbologyType(msg[1])
		d.Symbology.Mode = symbologyMode
		d.SyntaxType = BarcodeMessageScanData
		msg = msg[3:]
	} else {
		d.SyntaxType = BarcodeMessageFormat
	}

	subStr := msg[:]
	for len(subStr) > 0 {
		aiInfo, exists := detectAICode(subStr)

		if !exists {
			return Message{}, errors.New("error detecting valid AI code in message")
		}

		subStr = subStr[len(aiInfo.AI):]
		if aiInfo.IsFixedLength() {
			// Fixed length AI
			if len(subStr) >= aiInfo.Length() {
				value := subStr[:aiInfo.Length()]
				d.Elements = append(d.Elements, ElementString{
					aiInfo,
					value,
				})
				subStr = subStr[aiInfo.Length():]
			} else {
				return Message{}, errors.New("invalid GS1 identifier: insufficient length for AI " + aiInfo.AI)
			}
		} else {
			variableLengthAiEnd := strings.IndexRune(subStr, fnc1)
			if variableLengthAiEnd == -1 {
				variableLengthAiEnd = len(subStr)
			}

			value := subStr[:variableLengthAiEnd]
			d.Elements = append(d.Elements, ElementString{
				aiInfo,
				value,
			})

			if len(subStr) > len(value) {
				subStr = subStr[len(value)+1:]
			} else {
				subStr = subStr[len(value):]
			}
		}
	}

	return d, nil
}

func detectAICode(msg string) (ApplicationIdentifier, bool) {
	if len(msg) < 2 {
		return ApplicationIdentifier{}, false
	}
	if len(msg) > 4 {
		ai, ok := AIRegistry[msg[:4]]
		if ok {
			return ai, ok
		}
	}
	ai, ok := AIRegistry[msg[:2]]
	if ok {
		return ai, ok
	}
	return ApplicationIdentifier{}, false
}

// ParseElementString parses GS1 messages using the element string syntax. Example GS1 message compliant to this
// is `(01)09526064055028(17)250521(10)ABC123(21)456DEF`.
func ParseElementString(msg string) (d Message, _ error) {
	if !strings.HasPrefix(msg, "(") {
		return d, errors.New("invalid syntax: element string syntax must begin with '('")
	}

	subStr := msg[:]
	for len(subStr) > 0 {
		aiIDStart := strings.IndexRune(subStr, '(')
		aiIDEnd := strings.IndexRune(subStr, ')')

		aiID := subStr[aiIDStart+1 : aiIDEnd]
		ai, ok := AIRegistry[aiID]
		if !ok {
			return Message{}, fmt.Errorf("unkown AI: %s", aiID)
		}
		subStr = subStr[aiIDEnd+1:]

		// AI value's end
		aiDataEnd := strings.IndexRune(subStr, '(')
		if aiDataEnd == -1 {
			aiDataEnd = len(subStr)
		}

		aiData := subStr[:aiDataEnd]
		d.Elements = append(d.Elements, ElementString{
			ApplicationIdentifier: ai,
			DataField:             aiData,
		})

		subStr = subStr[aiDataEnd:]
	}

	if len(d.Elements) == 0 {
		return Message{}, errors.New("invalid syntax: no AI elements found")
	}

	d.SyntaxType = ElementStringSyntax
	return d, nil
}
