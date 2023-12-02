package models

type RoomStatusChanged struct {
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}
