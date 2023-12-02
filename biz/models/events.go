package models

import "time"

type RoomStatusChanged struct {
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}
