package domain

import "time"

type HiveState string

const (
	HiveColonized  HiveState = "colonized"
	HiveProduction HiveState = "production"
	HiveWintering  HiveState = "wintering"
	HiveRetired    HiveState = "retired"
)

type HiveModel string

const (
	ModelLangstroth HiveModel = "langstroth"
	ModelTopBar     HiveModel = "top_bar"
	ModelWarre      HiveModel = "warre"
)

// Hive 表示一只蜂箱及其当前蜂群状态。
type Hive struct {
	ID               string     `json:"id"`
	ApiaryID         string     `json:"apiaryId"`
	Model            HiveModel  `json:"model"`
	State            HiveState  `json:"state"`
	InstalledAt      time.Time  `json:"installedAt"`
	ActiveQueenID    string     `json:"activeQueenId"`
	LastInspectionAt *time.Time `json:"lastInspectionAt"`
}

// Transition 描述一次成功的蜂群状态迁移。
type Transition struct {
	HiveID string    `json:"hiveId"`
	From   HiveState `json:"from"`
	To     HiveState `json:"to"`
}

var hiveTransitions = map[HiveState]map[HiveState]bool{
	HiveColonized:  {HiveProduction: true, HiveRetired: true},
	HiveProduction: {HiveWintering: true, HiveRetired: true},
	HiveWintering:  {HiveProduction: true, HiveRetired: true},
	HiveRetired:    {},
}

func NewHive(id, apiaryID string, model HiveModel, installedAt time.Time) (Hive, error) {
	if err := ValidateIdentifier("id", id, 3, 40); err != nil {
		return Hive{}, err
	}
	if err := ValidateIdentifier("apiaryId", apiaryID, 3, 40); err != nil {
		return Hive{}, err
	}
	if !ValidHiveModel(model) {
		return Hive{}, Invalid("model", "model 必须是 %s、%s 或 %s", ModelLangstroth, ModelTopBar, ModelWarre)
	}
	return Hive{
		ID:          TrimText(id),
		ApiaryID:    TrimText(apiaryID),
		Model:       model,
		State:       HiveColonized,
		InstalledAt: installedAt.UTC(),
	}, nil
}

func ValidHiveModel(model HiveModel) bool {
	switch model {
	case ModelLangstroth, ModelTopBar, ModelWarre:
		return true
	default:
		return false
	}
}

func ValidHiveState(state HiveState) bool {
	switch state {
	case HiveColonized, HiveProduction, HiveWintering, HiveRetired:
		return true
	default:
		return false
	}
}

// CanTransitionHive 判断 from 到 to 是否为合法迁移，迁移到自身始终非法。
func CanTransitionHive(from, to HiveState) bool {
	targets, ok := hiveTransitions[from]
	if !ok {
		return false
	}
	return targets[to]
}

func (h Hive) AcceptsInspection() bool {
	return h.State != HiveRetired
}

func (h Hive) AcceptsQueenChange() bool {
	return h.State != HiveRetired
}
