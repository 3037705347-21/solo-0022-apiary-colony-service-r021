package apiary

import "example.com/apiary-colony-service/internal/domain"

// ParseHiveModel 把入口层字符串收敛为受支持的箱型。
func ParseHiveModel(value string) (domain.HiveModel, error) {
	model := domain.HiveModel(domain.TrimText(value))
	if !domain.ValidHiveModel(model) {
		return "", domain.Invalid(
			"model",
			"model 必须是 %s、%s 或 %s",
			domain.ModelLangstroth, domain.ModelTopBar, domain.ModelWarre,
		)
	}
	return model, nil
}

// ParseHiveState 把入口层字符串收敛为受支持的蜂群状态。
func ParseHiveState(value string) (domain.HiveState, error) {
	state := domain.HiveState(domain.TrimText(value))
	if !domain.ValidHiveState(state) {
		return "", domain.Invalid(
			"target",
			"target 必须是 %s、%s、%s 或 %s",
			domain.HiveColonized, domain.HiveProduction, domain.HiveWintering, domain.HiveRetired,
		)
	}
	return state, nil
}
