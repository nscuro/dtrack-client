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

func (c Client) UploadBOM(ctx context.Context, req BOMUploadRequest) (string, error) {
	return "", nil
}
