package models

import (
	"time"
)

type Shift struct {
	StudentID string     `json:"studentID"`
	Start     *time.Time `json:"start"`
	End       *time.Time `json:"end"`
}
