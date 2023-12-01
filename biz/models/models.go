package models

import "gorm.io/gorm"

var (
	RoomStatus_Empty    = "RoomStatus_Empty"
	RoomStatus_Occupied = "RoomStatus_Occupied"
)

type Room struct {
	gorm.Model

	Name   string `gorm:"name" db:"name"`
	Status string `gorm:"status" db:"status"`
}
