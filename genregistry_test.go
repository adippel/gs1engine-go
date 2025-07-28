package gs1

import (
	"bytes"
	"os"
	"slices"
	"testing"
)

func TestParseSyntaxDictionary(t *testing.T) {
	test := `
# AI     Flags Specification                    Attributes                                         Title
01         *?  N14,csum,keyoff1                 ex=255,37 dlpkey=22,10,21|235                    # GTIN
21             X..20                            req=01,03,8006 ex=235                            # SERIAL
410        *?  N13,csum,key                                                                      # SHIP TO LOC
8112        ?  X..70,couponposoffer
253         ?  N13,csum,key [X..17]             dlpkey                                           # GDTI
`
	ais, err := ParseSyntaxDictionary(bytes.NewBufferString(test))
	if err != nil {
		t.Error(err)
	}
	if len(ais) != 5 {
		t.Errorf("Expected 4 AIs, got %d", len(ais))
	}

	t.Run("ensure full AI description is parsed (AI-01)", func(t *testing.T) {
		ai01 := ais[0]
		if ai01.AI != "01" {
			t.Errorf("Expected AI 01, got %s", ai01.AI)
		}
		if ai01.Flags != "*?" {
			t.Errorf("Expected AI 01 flags to be *, got %s", ai01.Flags)
		}
		if !slices.Equal(ai01.Specification, []string{"N14", "csum", "keyoff1"}) {
			t.Errorf("Expected AI 01 specification to be N14,csum,keyoff1, got %s", ai01.Specification[0])
		}
		if !slices.Equal(ai01.Attributes, []string{"ex=255,37", "dlpkey=22,10,21|235"}) {
			t.Errorf("Expected AI 01 attributes to be ex=255,37, dlpkey=22,10,21|235, got %s", ai01.Attributes)
		}
		if ai01.Title != "GTIN" {
			t.Errorf("Expected AI 01 title to be GTIN, got %s", ai01.Title)
		}
	})

	t.Run("ensure AI description without FLAGS is parsed (AI-21)", func(t *testing.T) {
		ai21 := ais[1]
		if ai21.AI != "21" {
			t.Errorf("Expected AI 21, got %s", ai21.AI)
		}
		if ai21.Flags != "" {
			t.Errorf("Expected AI 21 flags to be EMPTY, got %s", ai21.Flags)
		}
		if !slices.Equal(ai21.Specification, []string{"X..20"}) {
			t.Errorf("Expected AI 21 specification to be X..20, got %s", ai21.Specification[0])
		}
		if !slices.Equal(ai21.Attributes, []string{"req=01,03,8006", "ex=235"}) {
			t.Errorf("Expected AI 21 attributes to be req=01,03,8006, ex=235, got %s", ai21.Attributes)
		}
		if ai21.Title != "SERIAL" {
			t.Errorf("Expected AI 21 title to be SERIAL, got %s", ai21.Title)
		}
	})

	t.Run("ensure AI description without attributes is parsed (AI-410)", func(t *testing.T) {
		ai410 := ais[2]
		if ai410.AI != "410" {
			t.Errorf("Expected AI 410, got %s", ai410.AI)
		}
		if ai410.Flags != "*?" {
			t.Errorf("Expected AI 410 flags to be *?, got %s", ai410.Flags)
		}
		if !slices.Equal(ai410.Specification, []string{"N13", "csum", "key"}) {
			t.Errorf("Expected AI 410 specification to be N13,csum,key, got %s", ai410.Specification[0])
		}
		if len(ai410.Attributes) != 0 {
			t.Errorf("Expected AI 410 attributes to be EMPTY, got %s", ai410.Attributes)
		}
		if ai410.Title != "SHIP TO LOC" {
			t.Errorf("Expected AI 410 title to be SHIP TO LOC, got %s", ai410.Title)
		}
	})

	t.Run("ensure AI description without title is parsed (AI-8112)", func(t *testing.T) {
		ai8112 := ais[3]
		if ai8112.AI != "8112" {
			t.Errorf("Expected AI 8112, got %s", ai8112.AI)
		}
		if ai8112.Flags != "?" {
			t.Errorf("Expected AI 8112 flags to be ?, got %s", ai8112.Flags)
		}
		if !slices.Equal(ai8112.Specification, []string{"X..70", "couponposoffer"}) {
			t.Errorf("Expected AI 8112 specification to be X..70, couponposoffer, got %s", ai8112.Specification[0])
		}
		if len(ai8112.Attributes) != 0 {
			t.Errorf("Expected AI 8112 attributes to be EMPTY, got %s", ai8112.Attributes)
		}
		if ai8112.Title != "" {
			t.Errorf("Expected AI 8112 title to be EMPTY, got %s", ai8112.Title)
		}
	})

	t.Run("ensure AI description with spec containing white space is parsed (AI-253)", func(t *testing.T) {
		ai253 := ais[4]
		if ai253.AI != "253" {
			t.Errorf("Expected AI 253, got %s", ai253.AI)
		}
		if ai253.Flags != "?" {
			t.Errorf("Expected AI 253 flags to be ?, got %s", ai253.Flags)
		}
		if !slices.Equal(ai253.Specification, []string{"N13", "csum", "key [X..17]"}) {
			t.Errorf("Expected AI253 specification to be N13,csum,key [X..17], got %s", ai253.Specification)
		}
		if !slices.Equal(ai253.Attributes, []string{"dlpkey"}) {
			t.Errorf("Expected AI 253 attributes to be dlpkey, got %s", ai253.Attributes)
		}
	})
}

func TestParseSyntaxDictionary_FullRegistryShouldBeBuild(t *testing.T) {
	data, err := os.Open("./testdata/gs1-syntax-dictionary-2025-01-30.txt")
	if err != nil {
		t.Error(err)
	}
	ais, err := ParseSyntaxDictionary(data)
	if err != nil {
		t.Error(err)
	}
	// 536 is retrieved from the tabular listing at https://ref.gs1.org/ai/
	if len(ais) != 536 {
		t.Errorf("Expected 536 AIs, got %d", len(ais))
	}
}
