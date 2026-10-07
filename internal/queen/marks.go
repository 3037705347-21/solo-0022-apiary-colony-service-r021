package queen

import "time"

const (
	MarkWhite  = "white"
	MarkYellow = "yellow"
	MarkRed    = "red"
	MarkGreen  = "green"
	MarkBlue   = "blue"
)

// ColorFor 按引入年份的个位返回约定的蜂王标记色。
func ColorFor(year int) string {
	switch ((year % 10) + 10) % 10 {
	case 1, 6:
		return MarkWhite
	case 2, 7:
		return MarkYellow
	case 3, 8:
		return MarkRed
	case 4, 9:
		return MarkGreen
	default:
		return MarkBlue
	}
}

func ColorForIntroduction(at time.Time) string {
	return ColorFor(at.Year())
}

func IsKnownColor(color string) bool {
	switch color {
	case MarkWhite, MarkYellow, MarkRed, MarkGreen, MarkBlue:
		return true
	default:
		return false
	}
}
