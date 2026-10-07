package status

import (
	"sort"
	"time"

	"example.com/apiary-colony-service/internal/domain"
	"example.com/apiary-colony-service/internal/store"
)

const InspectionWindow = 14 * 24 * time.Hour

type HiveSummary struct {
	HiveID           string     `json:"hiveId"`
	State            string     `json:"state"`
	LastInspectionAt *time.Time `json:"lastInspectionAt"`
	LatestConcerns   []string   `json:"latestConcerns"`
	QueenID          string     `json:"queenId"`
	QueenMarkColor   string     `json:"queenMarkColor"`
}

type Report struct {
	ApiaryID        string         `json:"apiaryId"`
	Name            string         `json:"name"`
	AsOf            time.Time      `json:"asOf"`
	HiveCount       int            `json:"hiveCount"`
	StateCounts     map[string]int `json:"stateCounts"`
	NeedsInspection []string       `json:"needsInspection"`
	ConcernHives    []string       `json:"concernHives"`
	LayingQueens    int            `json:"layingQueens"`
	Hives           []HiveSummary  `json:"hives"`
}

// BuildReport 在单份快照上完成聚合，保证计数与逐蜂箱摘要来自同一视图。
func BuildReport(snapshot store.Snapshot, apiary domain.Apiary, asOf time.Time) Report {
	asOf = asOf.UTC()
	report := Report{
		ApiaryID:        apiary.ID,
		Name:            apiary.Name,
		AsOf:            asOf,
		StateCounts:     initialStateCounts(),
		NeedsInspection: []string{},
		ConcernHives:    []string{},
		Hives:           []HiveSummary{},
	}
	hives := snapshot.ApiaryHives(apiary.ID)
	report.HiveCount = len(hives)
	for _, hive := range hives {
		fact, ok := snapshot.HiveFact(hive.ID)
		if !ok {
			continue
		}
		report.StateCounts[string(fact.State)]++
		summary := HiveSummary{
			HiveID:           hive.ID,
			State:            string(fact.State),
			LastInspectionAt: fact.LastInspectionAt,
			LatestConcerns:   []string{},
		}
		if fact.ActiveQueenID != "" {
			summary.QueenID = fact.ActiveQueenID
			summary.QueenMarkColor = fact.QueenMarkColor
			report.LayingQueens++
		}
		if len(fact.LatestConcerns) > 0 {
			summary.LatestConcerns = append([]string{}, fact.LatestConcerns...)
			report.ConcernHives = append(report.ConcernHives, hive.ID)
		}
		if NeedsInspectionAt(fact.LastInspectionAt, asOf) {
			report.NeedsInspection = append(report.NeedsInspection, hive.ID)
		}
		report.Hives = append(report.Hives, summary)
	}
	sort.Strings(report.NeedsInspection)
	sort.Strings(report.ConcernHives)
	return report
}

// NeedsInspection 判定蜂箱是否需要巡检：从未巡检，或距 asOf 超过巡检窗口。
func NeedsInspection(hive domain.Hive, asOf time.Time) bool {
	return NeedsInspectionAt(hive.LastInspectionAt, asOf)
}

// NeedsInspectionAt 判定最近一次巡检时间是否需要巡检：从未巡检，或距 asOf 超过巡检窗口。
func NeedsInspectionAt(lastInspectionAt *time.Time, asOf time.Time) bool {
	if lastInspectionAt == nil {
		return true
	}
	return asOf.Sub(lastInspectionAt.UTC()) > InspectionWindow
}

func initialStateCounts() map[string]int {
	return map[string]int{
		string(domain.HiveColonized):  0,
		string(domain.HiveProduction): 0,
		string(domain.HiveWintering):  0,
		string(domain.HiveRetired):    0,
	}
}
