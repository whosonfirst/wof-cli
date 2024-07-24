package uris

import (
	"testing"
)

func TestMatchWhosOnFirstId(t *testing.T) {

	tests := []string{
		"1",
		"1-alt-nixta",
	}

	for _, str := range tests {

		if !re_wofid.MatchString(str) {
			t.Fatalf("Failed to match '%s'", str)
		}
	}

}
