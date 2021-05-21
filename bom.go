package dtrack

import (
	"context"

	"github.com/google/uuid"
)

type BOMUploadRequest struct {
	ProjectUUID    *uuid.UUID `json:"project,omitempty"`
	ProjectName    string     `json:"projectName,omitempty"`
	ProjectVersion string     `json:"projectVersion,omitempty"`
	AutoCreate     bool       `json:"autoCreate"`
	BOM            string     `json:"bom"`
}

type BOMService interface {
	ExportProjectAsCycloneDX(ctx context.Context, pUUID uuid.UUID) (string, error)
	IsProcessingBOM(ctx context.Context, token string) (bool, error)
	UploadBOM(ctx context.Context, req BOMUploadRequest) (string, error)
}
