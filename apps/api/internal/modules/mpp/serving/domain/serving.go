package domain

import "time"

// Completion is the result of finishing a service. Duration is computed from the
// antrian timestamps (served_at → done_at); neither table stores it.
type Completion struct {
	AntrianID       string    `json:"antrian_id"`
	Nomor           string    `json:"nomor"`
	Status          string    `json:"status"`
	DoneAt          time.Time `json:"done_at"`
	DurationSeconds *int      `json:"duration_seconds"`
}

// CurrentCall is the big number on the TV, with the loket to walk to.
type CurrentCall struct {
	AntrianID string    `json:"antrian_id"`
	Nomor     string    `json:"nomor"`
	Loket     string    `json:"loket"`
	Status    string    `json:"status"`
	CallCount int       `json:"call_count"`
	TTSText   string    `json:"tts_text"`
	CalledAt  time.Time `json:"called_at"`
}

// NextUp is one upcoming number in the list under the current call.
type NextUp struct {
	AntrianID string `json:"antrian_id"`
	Nomor     string `json:"nomor"`
}

// Display is one TV snapshot. Current is nil between calls — an idle screen, not
// an error (FR-TV-01).
type Display struct {
	InstansiID string       `json:"instansi_id"`
	Current    *CurrentCall `json:"current"`
	NextUp     []NextUp     `json:"next_up"`
}
