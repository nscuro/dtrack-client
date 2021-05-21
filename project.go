package dtrack

import (
	"context"

	"github.com/google/uuid"
)

type Project struct {
	UUID    uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Group   string    `json:"group"`
}

type ProjectService interface {
	GetAllProjects(ctx context.Context) ([]Project, error)
	GetProjectByUUID(ctx context.Context, pUUID uuid.UUID) (*Project, error)
	LookupProject(ctx context.Context, name, version string) (*Project, error)
}
