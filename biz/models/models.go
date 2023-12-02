package models

var (
	RoomStatus_Empty    = "EMPTY"
	RoomStatus_Occupied = "OCCUPIED"
)

type Room struct {
	Common

	Name   string `gorm:"name" db:"name" json:"name"`
	Status string `gorm:"status" db:"status" json:"status"`
}

type Common struct {
	Id        uint32 `gorm:"id,primaryKey" db:"id" json:"id"`
	CreatedAt string `gorm:"created_at" db:"created_at" json:"created_at"`
	UpdatedAt string `gorm:"updated_at" db:"updated_at" json:"updated_at"`
}
