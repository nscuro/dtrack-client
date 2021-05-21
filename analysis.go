package dtrack

import (
	"context"

	"github.com/google/uuid"
)

type Analysis struct {
	Comments   []AnalysisComment `json:"comments"`
	State      string            `json:"analysisState"`
	Suppressed bool              `json:"isSuppressed"`
}

type AnalysisComment struct {
	Comment   string `json:"comment"`
	Commenter string `json:"commenter"`
	Timestamp string `json:"timestamp"`
}

type AnalysisRequest struct {
	ComponentUUID     uuid.UUID `json:"component"`
	ProjectUUID       uuid.UUID `json:"project"`
	VulnerabilityUUID uuid.UUID `json:"vulnerability"`
	Comment           string    `json:"comment,omitempty"`
	State             string    `json:"analysisState,omitempty"`
	Suppressed        *bool     `json:"isSuppressed,omitempty"`
}

type AnalysisService interface {
	CreateAnalysis(ctx context.Context, req AnalysisRequest) (*Analysis, error)
	GetAnalysis(ctx context.Context, cUUID, pUUID, vUUID uuid.UUID) (*Analysis, error)
}
