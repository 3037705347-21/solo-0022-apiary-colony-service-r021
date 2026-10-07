package status

import (
	"time"

	"example.com/apiary-colony-service/internal/domain"
	"example.com/apiary-colony-service/internal/store"
)

// Service 读取单份快照并给出蜂场状态汇总。
type Service struct {
	store *store.Memory
	clock func() time.Time
}

func New(memory *store.Memory, clock func() time.Time) *Service {
	if clock == nil {
		clock = time.Now
	}
	return &Service{store: memory, clock: clock}
}

func (s *Service) ApiaryStatus(apiaryID string, asOf *time.Time) (Report, error) {
	snapshot := s.store.Snapshot()
	apiary, ok := snapshot.Apiaries[apiaryID]
	if !ok {
		return Report{}, domain.NotFound("蜂场", apiaryID)
	}
	moment := s.clock()
	if asOf != nil {
		moment = *asOf
	}
	return BuildReport(snapshot, apiary, moment), nil
}
