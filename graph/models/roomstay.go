package models

import "time"

type RoomStay struct {
	StayID    int     `json:"stayID" gorm:"type:serial auto_increment;not null;primary_key"`
	Room      string     `json:"room"`
	StudentID string     `json:"studentID"`
	GroupSize *int       `json:"groupSize"`
	Start     *time.Time `json:"start"`
	End       *time.Time `json:"end"`
}