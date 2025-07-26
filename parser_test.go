package gs1

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseElementStringSyntax(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name            string
		args            args
		wantDataMessage DataMessage
		wantErr         bool
	}{
		{
			name: "Compliant syntax with single AI SHOULD parse",
			args: args{"(01)09526064055028"},
			wantDataMessage: DataMessage{
				SyntaxType: ElementStringSyntax,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "09526064055028"},
				},
			},
		},
		{
			name: "Compliant syntax with multiple AIs SHOULD parse",
			args: args{"(01)09526064055028(21)123456"},
			wantDataMessage: DataMessage{
				SyntaxType: ElementStringSyntax,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "09526064055028"},
					{ApplicationIdentifierInfo{AI21}, "123456"},
				},
			},
			wantErr: false,
		},
		{
			name:    "Missing opening parenthesis at start SHOULD return an error",
			args:    args{"123456"},
			wantErr: true,
		},
		{
			name:    "Usage of unknown AI SHOULD return an error",
			args:    args{"(0000)123456"},
			wantErr: true,
		},
		{
			name:    "Empty AI code parenthesis SHOULD return an error",
			args:    args{"(01)123456()66666"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := ParseElementString(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseElementString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantDataMessage) {
				t.Errorf("ParseElementString() gotD = %v, want %v", gotD, tt.wantDataMessage)
			}
		})
	}
}

func TestParseBarcodeMessage(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantD   DataMessage
		wantErr bool
	}{
		{
			name: "Message with {GS} and ^ characters to indicate FNC1 SHOULD parse",
			args: args{"^010952606405502810654321{GS}21123456"},
			wantD: DataMessage{
				SyntaxType: BarcodeMessageFormat,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "09526064055028"},
					{ApplicationIdentifierInfo{AI10}, "654321"},
					{ApplicationIdentifierInfo{AI21}, "123456"},
				},
			},
		},
		{
			name: "Message with a single element string SHOULD parse",
			args: args{"^0109526064055028"},
			wantD: DataMessage{
				SyntaxType: BarcodeMessageFormat,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "09526064055028"},
				},
			},
			wantErr: false,
		},
		{
			name: "Message with a single element string with fixed and variable length AIs SHOULD parse",
			args: args{"^010952606405502821123456"},
			wantD: DataMessage{
				SyntaxType: BarcodeMessageFormat,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "09526064055028"},
					{ApplicationIdentifierInfo{AI21}, "123456"},
				},
			},
			wantErr: false,
		},
		{
			name: "Message with symbology information SHOULD parse",
			args: args{fmt.Sprintf("%c%s%d", SymbologyFlag, GS1DataMatrix, 2) + "0109526064055028"},
			wantD: DataMessage{
				SyntaxType: BarcodeMessageScanData,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "09526064055028"},
				},
				Symbology: SymbologyIdentifier{
					Type: GS1DataMatrix,
					Mode: 2,
				},
			},
		},
		{
			name: "Messages with 4-digit AIs SHOULD parse",
			args: args{"356412345621123456"},
			wantD: DataMessage{
				SyntaxType: BarcodeMessageFormat,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI3564}, "123456"},
					{ApplicationIdentifierInfo{AI21}, "123456"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := ParseBarcodeMessage(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBarcodeMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("ParseBarcodeMessage() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func TestParseDataMessage(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantD   DataMessage
		wantErr bool
	}{
		{
			name: "Element String syntax SHOULD parse",
			args: args{"(01)76100120000010(10)10002256"},
			wantD: DataMessage{
				SyntaxType: ElementStringSyntax,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "76100120000010"},
					{ApplicationIdentifierInfo{AI10}, "10002256"},
				},
			},
		},
		{
			name: "Barcode message syntax SHOULD parse",
			args: args{"^0109526064055028"},
			wantD: DataMessage{
				SyntaxType: BarcodeMessageFormat,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "09526064055028"},
				},
			},
		},
		{
			name: "Barcode message scan data syntax SHOULD parse",
			args: args{"]C1010123456789012815057072"},
			wantD: DataMessage{
				SyntaxType: BarcodeMessageScanData,
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "01234567890128"},
					{ApplicationIdentifierInfo{AI15}, "057072"},
				},
				Symbology: SymbologyIdentifier{
					Type: GS1128,
					Mode: 1,
				},
			},
		},
		{
			name:    "Digital-Link URI Syntax SHOULD return error",
			args:    args{"https://example.com/01/09526064055028"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := ParseMessage(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("ParseMessage() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}
