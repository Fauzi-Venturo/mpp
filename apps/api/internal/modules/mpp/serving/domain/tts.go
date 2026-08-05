package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// The backend builds tts_text so every TV pronounces a number the same way —
// docs/05-integrations/tts-voice-calling.md:54. Digits are spoken one by one, the
// instansi prefix stays a letter, and the loket becomes a number word.
//
// ponytail: the phrasing is hardcoded. FR-CFG-03 wants it read from
// mpp.system_config (config_key 'tts_text'), but that row is not seeded yet.
var digitWords = [10]string{"nol", "satu", "dua", "tiga", "empat", "lima", "enam", "tujuh", "delapan", "sembilan"}

// numberWords covers the loket labels a building realistically has; anything
// larger falls back to digit-by-digit, which is still understandable.
var numberWords = map[int]string{
	1: "satu", 2: "dua", 3: "tiga", 4: "empat", 5: "lima",
	6: "enam", 7: "tujuh", 8: "delapan", 9: "sembilan", 10: "sepuluh",
	11: "sebelas", 12: "dua belas",
}

// BuildTTSText renders "A-014" at "Loket 3" as
// "Nomor antrian A - nol satu empat, silakan menuju loket tiga".
func BuildTTSText(nomor, loket string) string {
	prefix, digits, _ := strings.Cut(nomor, "-")

	spoken := make([]string, 0, len(digits))
	for _, r := range digits {
		if r >= '0' && r <= '9' {
			spoken = append(spoken, digitWords[r-'0'])
		}
	}

	text := fmt.Sprintf("Nomor antrian %s - %s", prefix, strings.Join(spoken, " "))
	if label := spokenLoket(loket); label != "" {
		text += ", silakan menuju loket " + label
	}

	return text
}

// spokenLoket turns "Loket 3" into "tiga". A label without a number is spoken as
// written, so an unusual name still reaches the citizen.
func spokenLoket(loket string) string {
	fields := strings.Fields(loket)
	if len(fields) == 0 {
		return ""
	}

	last := fields[len(fields)-1]
	n, err := strconv.Atoi(last)
	if err != nil {
		return strings.ToLower(loket)
	}

	if word, ok := numberWords[n]; ok {
		return word
	}

	spoken := make([]string, 0, len(last))
	for _, r := range last {
		spoken = append(spoken, digitWords[r-'0'])
	}
	return strings.Join(spoken, " ")
}
