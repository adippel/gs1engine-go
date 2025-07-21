package gs1

import "strings"

type ApplicationIdentifier struct {
	AIDescription
}

// RequiresFNC1Separator checks if the AI description declares the flag '*'.
func (ai ApplicationIdentifier) RequiresFNC1Separator() bool {
	return strings.ContainsRune(ai.Flags, '*')
}

func (ai ApplicationIdentifier) HasFixedLength() bool {
	return ai.RequiresFNC1Separator()
}
