package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"example.com/apiary-colony-service/internal/domain"
)

const maxBodyBytes = 1 << 20

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	body := http.MaxBytesReader(w, r.Body, maxBodyBytes)
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		if errors.Is(err, io.EOF) {
			return domain.Invalid("body", "请求体不能为空")
		}
		return domain.Invalid("body", "请求体不是合法 JSON：%v", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return domain.Invalid("body", "请求体只能包含一个 JSON 对象")
	}
	return nil
}

// parseOptionalTime 解析可缺省的时间参数，缺省时返回 nil 交由调用方取当前时间。
func parseOptionalTime(field, value string) (*time.Time, error) {
	trimmed := domain.TrimText(value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return nil, domain.Invalid(field, "%s 必须是 RFC3339 时间", field)
	}
	utc := parsed.UTC()
	return &utc, nil
}
