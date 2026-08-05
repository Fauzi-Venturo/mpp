package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	goredis "github.com/redis/go-redis/v9"

	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/domain"
	"github.com/ndollem/mpp/apps/api/internal/modules/mpp/queue/repository"
)

// counterTTL outlives one operating day so the daily reset comes from the date in
// the key (BR-03), not from expiry timing.
const counterTTL = 48 * time.Hour

// nomorFormat renders A-014. ponytail: BR-04 says the pattern is configurable, but
// mpp.system_config has no key for it yet — read it from there once it exists.
const nomorFormat = "%s-%03d"

type QueueService struct {
	antrianRepo *repository.AntrianRepository
	rdb         *goredis.Client
}

func NewQueueService(antrianRepo *repository.AntrianRepository, rdb *goredis.Client) *QueueService {
	return &QueueService{antrianRepo: antrianRepo, rdb: rdb}
}

// seqKey scopes the counter per service per day: one stream per jenis layanan
// (BR-02), starting over at midnight (BR-03).
func seqKey(layananID string, day time.Time) string {
	return fmt.Sprintf("queue:seq:%s:%s", layananID, day.Format("2006-01-02"))
}

// EnqueueTx allocates a number and writes the ticket inside the caller's
// transaction, so a caller that also mutates a booking can roll both back.
//
// The number comes from Redis INCR, which is what makes concurrent allocation
// safe — no read-then-write anywhere. SetNX seeds a cold counter from the highest
// number Postgres already holds, so a restarted Redis (or a day with seeded demo
// rows) resumes instead of colliding.
//
// ponytail: a rolled-back transaction burns its number, leaving a gap in the
// sequence. Gaps are harmless; duplicates are not. Reclaim them only if the
// sequence ever has to be gapless.
func (s *QueueService) EnqueueTx(ctx context.Context, tx pgx.Tx, p domain.EnqueueParams) (*domain.Antrian, error) {
	prefix, maxSeq, err := s.antrianRepo.PrefixAndMaxSeqTx(ctx, tx, p.InstansiID, p.JenisLayananID, p.QueueDate)
	if err != nil {
		return nil, err
	}

	key := seqKey(p.JenisLayananID, p.QueueDate)
	if err := s.rdb.SetNX(ctx, key, maxSeq, counterTTL).Err(); err != nil {
		return nil, fmt.Errorf("seed queue counter: %w", err)
	}

	seq, err := s.rdb.Incr(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("allocate queue number: %w", err)
	}

	if p.Source == "" {
		p.Source = domain.SourceBooking
	}

	return s.antrianRepo.InsertTx(ctx, tx, p, fmt.Sprintf(nomorFormat, prefix, seq), int(seq))
}

// ListWaiting returns today's waiting stream for one service.
func (s *QueueService) ListWaiting(ctx context.Context, layananID string) ([]domain.Antrian, error) {
	return s.antrianRepo.ListWaiting(ctx, layananID)
}
