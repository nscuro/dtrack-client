package dtrack

import (
	"github.com/google/uuid"
)

type Component struct {
	UUID    uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Group   string    `json:"group"`
}
