package gs1

import "testing"

func TestApplicationIdentifierInfo_RequiresFNC1Separator(t *testing.T) {
	type fields struct {
		AIDescription AIDescription
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "AI with fixed length requires no FNC1 terminator",
			fields: fields{AI01},
			want:   false,
		},
		{
			name:   "AI with non * flag is detected as variable length",
			fields: fields{AI10},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := ApplicationIdentifierInfo{
				AIDescription: tt.fields.AIDescription,
			}
			if got := ai.RequiresFNC1Separator(); got != tt.want {
				t.Errorf("RequiresFNC1Separator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplicationIdentifierInfo_Length(t *testing.T) {
	type fields struct {
		AIDescription AIDescription
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "AI with fixed length return valid length",
			fields: fields{AIDescription: AI01},
			want:   14,
		},
		{
			name:   "AI with variable length returns -1",
			fields: fields{AIDescription: AI10},
			want:   -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai := ApplicationIdentifierInfo{
				AIDescription: tt.fields.AIDescription,
			}
			if got := ai.Length(); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataMessage_AsElementString(t *testing.T) {
	type fields struct {
		Symbology  SymbologyIdentifier
		SyntaxType DataSyntaxType
		Elements   []ElementString
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Valid GS1 data with two AIs SHOULD return valid element string",
			fields: fields{
				Elements: []ElementString{
					{ApplicationIdentifierInfo{AI01}, "01234567890128"},
					{ApplicationIdentifierInfo{AI15}, "057072"},
				},
			},
			want: "(01)01234567890128(15)057072",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DataMessage{
				Symbology:  tt.fields.Symbology,
				SyntaxType: tt.fields.SyntaxType,
				Elements:   tt.fields.Elements,
			}
			if got := d.AsElementString(); got != tt.want {
				t.Errorf("AsElementString() = %v, want %v", got, tt.want)
			}
		})
	}
}
