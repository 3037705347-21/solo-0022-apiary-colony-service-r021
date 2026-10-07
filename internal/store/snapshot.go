package store

import (
	"sort"

	"example.com/apiary-colony-service/internal/domain"
)

// Snapshot 是一次一致性读取的结果，供聚合与校验使用。
type Snapshot struct {
	Apiaries          map[string]domain.Apiary
	Hives             map[string]domain.Hive
	InspectionsByHive map[string][]domain.Inspection
	QueensByHive      map[string][]domain.Queen
	ReadFacts         map[string]HiveReadFact
}

// HiveFact 返回某只蜂箱的派生读取事实；快照没有携带时按快照内的权威记录装配。
func (s Snapshot) HiveFact(hiveID string) (HiveReadFact, bool) {
	if fact, ok := s.ReadFacts[hiveID]; ok {
		return cloneReadFact(fact), true
	}
	hive, ok := s.Hives[hiveID]
	if !ok {
		return HiveReadFact{}, false
	}
	fact := HiveReadFact{
		HiveID:           hive.ID,
		State:            hive.State,
		ActiveQueenID:    hive.ActiveQueenID,
		LastInspectionAt: cloneTime(hive.LastInspectionAt),
		LatestConcerns:   []string{},
	}
	if queen, ok := activeQueen(s.QueensByHive[hiveID], hive.ActiveQueenID); ok {
		fact.QueenMarkColor = queen.MarkColor
	}
	if latest, ok := latestInspection(s.InspectionsByHive[hiveID]); ok {
		observedAt := latest.ObservedAt
		fact.LastInspectionAt = &observedAt
		fact.LatestConcerns = append([]string{}, latest.Concerns...)
	}
	return fact, true
}

// ApiaryHives 返回蜂场下的蜂箱，按编号升序稳定排序。
func (s Snapshot) ApiaryHives(apiaryID string) []domain.Hive {
	hives := make([]domain.Hive, 0, len(s.Hives))
	for _, hive := range s.Hives {
		if hive.ApiaryID == apiaryID {
			hives = append(hives, hive)
		}
	}
	sort.Slice(hives, func(i, j int) bool { return hives[i].ID < hives[j].ID })
	return hives
}

// LatestInspection 返回蜂箱最近一次巡检；观察时间相同时按编号取较大者。
func (s Snapshot) LatestInspection(hiveID string) (domain.Inspection, bool) {
	return latestInspection(s.InspectionsByHive[hiveID])
}

// LayingQueen 返回蜂箱当前处于产卵状态的蜂王。
func (s Snapshot) LayingQueen(hiveID string) (domain.Queen, bool) {
	for _, item := range s.QueensByHive[hiveID] {
		if item.IsLaying() {
			return item, true
		}
	}
	return domain.Queen{}, false
}

func (s Snapshot) Queens(hiveID string) []domain.Queen {
	items := s.QueensByHive[hiveID]
	out := make([]domain.Queen, len(items))
	copy(out, items)
	return out
}
