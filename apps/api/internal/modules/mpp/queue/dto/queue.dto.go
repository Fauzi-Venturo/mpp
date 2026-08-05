package dto

import (
	"time"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/domain"
)

// QueueQuery selects the stream to read. One service = one stream (BR-02).
type QueueQuery struct {
	LayananID string `form:"layanan_id" binding:"required,uuid"`
}

// QueueItem is one ticket as the loket/TV clients see it. Field names follow the
// websocket contract (docs/04-api/websocket-events.md) so slice 6 does not have to
// rename anything.
type QueueItem struct {
	AntrianID string    `json:"antrian_id"`
	Nomor     string    `json:"nomor"`
	Status    string    `json:"status"`
	QueuedAt  time.Time `json:"queued_at"`
}

// QueueStream mirrors the `snapshot` websocket payload: the waiting list plus its
// count, keyed by service. The shared response envelope's meta only carries
// pagination, so the count travels in data.
type QueueStream struct {
	LayananID    string      `json:"layanan_id"`
	WaitingCount int         `json:"waiting_count"`
	Waiting      []QueueItem `json:"waiting"`
}

// NewQueueStream projects domain rows onto the wire shape.
func NewQueueStream(layananID string, rows []domain.Antrian) QueueStream {
	items := make([]QueueItem, 0, len(rows))
	for _, a := range rows {
		items = append(items, QueueItem{
			AntrianID: a.ID,
			Nomor:     a.Nomor,
			Status:    string(a.Status),
			QueuedAt:  a.QueuedAt,
		})
	}

	return QueueStream{LayananID: layananID, WaitingCount: len(items), Waiting: items}
}
