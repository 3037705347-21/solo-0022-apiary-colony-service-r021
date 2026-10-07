package domain

import "time"

type QueenState string

const (
	QueenLaying     QueenState = "laying"
	QueenSuperseded QueenState = "superseded"
	QueenLost       QueenState = "lost"
)

type QueenOrigin string

const (
	OriginBred  QueenOrigin = "bred"
	OriginCells QueenOrigin = "cells"
	OriginSwarm QueenOrigin = "swarm"
)

// Queen 是蜂箱内的一只蜂王记录，关闭后不可再改。
type Queen struct {
	ID           string      `json:"id"`
	HiveID       string      `json:"hiveId"`
	IntroducedAt time.Time   `json:"introducedAt"`
	Origin       QueenOrigin `json:"origin"`
	MarkColor    string      `json:"markColor"`
	State        QueenState  `json:"state"`
	ClosedAt     *time.Time  `json:"closedAt"`
}

func NewQueen(hiveID string, introducedAt time.Time, origin QueenOrigin, markColor string) (Queen, error) {
	if err := ValidateIdentifier("hiveId", hiveID, 3, 40); err != nil {
		return Queen{}, err
	}
	if introducedAt.IsZero() {
		return Queen{}, Invalid("introducedAt", "introducedAt 不能为空")
	}
	if !ValidQueenOrigin(origin) {
		return Queen{}, Invalid("origin", "origin 必须是 %s、%s 或 %s", OriginBred, OriginCells, OriginSwarm)
	}
	if err := ValidateText("markColor", markColor, 1, 20); err != nil {
		return Queen{}, err
	}
	return Queen{
		HiveID:       TrimText(hiveID),
		IntroducedAt: introducedAt.UTC(),
		Origin:       origin,
		MarkColor:    TrimText(markColor),
		State:        QueenLaying,
	}, nil
}

func ValidQueenOrigin(origin QueenOrigin) bool {
	switch origin {
	case OriginBred, OriginCells, OriginSwarm:
		return true
	default:
		return false
	}
}

func (q Queen) IsLaying() bool {
	return q.State == QueenLaying
}
