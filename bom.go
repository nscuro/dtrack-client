package dtrack

import (
	"context"
	"fmt"
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

type BOMService struct {
	client *Client
}

func (bs BOMService) Upload(ctx context.Context, uploadReq BOMUploadRequest) (BOMUploadToken, error) {
	req, err := bs.client.newRequest(ctx, http.MethodPut, "/api/v1/bom", withBody(uploadReq))
	if err != nil {
		return "", err
	}

	var uploadRes bomUploadResponse
	_, err = bs.client.doRequest(req, &uploadRes)
	if err != nil {
		return "", err
	}

	return uploadRes.Token, nil
}

type bomProcessingResponse struct {
	Processing bool `json:"processing"`
}

func (bs BOMService) IsBeingProcessed(ctx context.Context, token BOMUploadToken) (bool, error) {
	req, err := bs.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/bom/token/%s", token))
	if err != nil {
		return false, err
	}

	var processingResponse bomProcessingResponse
	_, err = bs.client.doRequest(req, &processingResponse)
	if err != nil {
		return false, err
	}

	return processingResponse.Processing, nil
}
