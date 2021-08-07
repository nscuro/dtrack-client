package dtrack

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type BOMUploadRequest struct {
	ProjectUUID    *uuid.UUID `json:"project,omitempty"`
	ProjectName    string     `json:"projectName,omitempty"`
	ProjectVersion string     `json:"projectVersion,omitempty"`
	AutoCreate     bool       `json:"autoCreate"`
	BOM            string     `json:"bom"`
}

type bomUploadResponse struct {
	Token BOMUploadToken `json:"token"`
}

type BOMUploadToken string

func (c Client) UploadBOM(ctx context.Context, uploadReq BOMUploadRequest) (BOMUploadToken, error) {
	req, err := c.newRequest(ctx, http.MethodPut, "/api/v1/bom", withBody(uploadReq))
	if err != nil {
		return "", err
	}

	var uploadRes bomUploadResponse
	_, err = c.doRequest(req, &uploadRes)
	if err != nil {
		return "", err
	}

	return uploadRes.Token, nil
}
