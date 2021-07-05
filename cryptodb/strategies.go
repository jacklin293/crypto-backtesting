package cryptodb

import "time"

type Strategy struct {
	Id        int64
	Title     string
	Desc      string // Description
	Type      string // strategy type
	Params    string // JSON string to store ad-hoc params of strategy type
	Enabled   int    // 0: disabled 1: enabled
	CreatedAt time.Time
	UpdatedAt time.Time
}
