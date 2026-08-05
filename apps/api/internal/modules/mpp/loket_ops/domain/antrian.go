package domain

import "time"

// AntrianStatus mirrors the mpp.antrian_status enum. Only the values slice 5
// transitions between are declared here; the queue module owns its own subset.
type AntrianStatus string

const (
	StatusWaiting AntrianStatus = "WAITING"
	StatusCalled  AntrianStatus = "CALLED"
	StatusServing AntrianStatus = "SERVING"
	StatusSkipped AntrianStatus = "SKIPPED"
)

// MaxCallCount is the call-3x-then-skip rule (BR-16 / FR-OPR-03). The database
// carries the same limit as a CHECK constraint; this guard keeps a 4th call a
// clean 409 instead of a constraint violation.
const MaxCallCount = 3

// Called is one antrian as the operator app sees it after a call action.
type Called struct {
	AntrianID string        `json:"antrian_id"`
	Nomor     string        `json:"nomor"`
	Status    AntrianStatus `json:"status"`
	LoketID   string        `json:"loket_id"`
	Loket     string        `json:"loket"`
	CallCount int           `json:"call_count"`
	CalledAt  *time.Time    `json:"called_at,omitempty"`
}
