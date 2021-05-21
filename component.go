package dtrack

import (
	"context"

	"github.com/google/uuid"
)

type Component struct {
	UUID    uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Group   string    `json:"group"`
}

type ComponentService interface {
	GetComponentsForProject(ctx context.Context, pUUID uuid.UUID) ([]Component, error)
	GetComponentByUUID(ctx context.Context, cUUID uuid.UUID) (*Component, error)
	GetComponentByHash(ctx context.Context, hash string) (*Component, error)
}
