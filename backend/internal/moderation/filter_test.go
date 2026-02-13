package moderation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsProfanity_KnownWord(t *testing.T) {
	assert.True(t, ContainsProfanity("this is shit"), "should detect known profanity word")
	assert.True(t, ContainsProfanity("what the fuck"), "should detect known profanity word")
	assert.True(t, ContainsProfanity("DAMN it"), "should be case-insensitive")
}

func TestContainsProfanity_SubstringNotMatched(t *testing.T) {
	// "class" contains "ass" as a substring but is a different word.
	// The filter uses whole-word matching via tokenization, so it should NOT flag "class".
	assert.False(t, ContainsProfanity("class"), "should not match 'ass' inside 'class'")
	assert.False(t, ContainsProfanity("grasshopper"), "should not match 'ass' inside 'grasshopper'")
	assert.False(t, ContainsProfanity("classic art"), "should not match 'ass' inside 'classic'")
	assert.False(t, ContainsProfanity("assume nothing"), "should not match 'ass' inside 'assume'")
	assert.False(t, ContainsProfanity("cockatoo"), "should not match 'cock' inside 'cockatoo'")
}

func TestContainsProfanity_CleanText(t *testing.T) {
	assert.False(t, ContainsProfanity("Hello, this is a perfectly clean sentence."))
	assert.False(t, ContainsProfanity("Great work on the project!"))
	assert.False(t, ContainsProfanity(""))
}

func TestContainsProfanity_MultiWordPhrase(t *testing.T) {
	assert.True(t, ContainsProfanity("you son of a bitch"), "should detect multi-word profanity phrase")
}

func TestScanContent_NoProfanity(t *testing.T) {
	reasons := ScanContent("Beautiful Sunset", "A gorgeous sunset over the ocean.")
	assert.Empty(t, reasons, "clean content should return no flags")
}

func TestScanContent_ProfanityInTitle(t *testing.T) {
	reasons := ScanContent("What the hell", "A perfectly normal description.")
	assert.Contains(t, reasons, "profanity_in_title")
	assert.NotContains(t, reasons, "profanity_in_description")
}

func TestScanContent_ProfanityInDescription(t *testing.T) {
	reasons := ScanContent("Nice Title", "This is total bullshit")
	assert.NotContains(t, reasons, "profanity_in_title")
	assert.Contains(t, reasons, "profanity_in_description")
}

func TestScanContent_ProfanityInBoth(t *testing.T) {
	reasons := ScanContent("Damn photo", "This is total crap")
	assert.Contains(t, reasons, "profanity_in_title")
	assert.Contains(t, reasons, "profanity_in_description")
	assert.Len(t, reasons, 2)
}
