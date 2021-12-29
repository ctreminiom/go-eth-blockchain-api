package models

type Error struct {
	Code    uint64 `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
