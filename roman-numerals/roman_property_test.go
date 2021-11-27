package roman

import (
	"testing"
	"log"
	"testing/quick"
)


func TestPropertiesOfConversion(t *testing.T) {
	assertion := func(arabic uint16) bool {
		if arabic > 3999 {
			return true
		}
		log.Println("testing", arabic)
		roman := ConvertToRoman(int(arabic))
		res := ConvertToNumeral(roman)

		return int(arabic) == res
	}

	if err := quick.Check(assertion, &quick.Config{MaxCount: 300}); err != nil {
		t.Error("failed checks", err)
	}
}