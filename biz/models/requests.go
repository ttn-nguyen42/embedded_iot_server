package models

type ListCommon struct {
	Page  uint64 `json:"-"`
	Limit uint64 `json:"-"`
}

type GetRoomsRequest struct {
	ListCommon
}

type GetRoomsResponse struct {
	Rooms []Room `json:"rooms"`
}

type GetRoomRequest struct {
	Id uint32 `json:"-"`
}

type GetRoomResponse struct {
	Room *Room `json:"room"`
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type UpdateRoomRequest struct {
	Id     uint32 `json:"-"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CreateRoomResponse struct {
	Id uint32 `json:"id"`
}

type UpdateRoomResponse struct {
	Id uint32 `json:"id"`
}
