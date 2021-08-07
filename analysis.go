package dtrack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type AnalysisState int

const (
	AnalysisStateNotSet AnalysisState = iota
	AnalysisStateInTriage
	AnalysisStateExploitable
	AnalysisStateNotAffected
	AnalysisStateFalsePositive
)

func (a AnalysisState) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *AnalysisState) UnmarshalJSON(bytes []byte) error {
	var val string
	if err := json.Unmarshal(bytes, &val); err != nil {
		return err
	}

	switch val {
	case AnalysisStateNotSet.String():
		*a = AnalysisStateNotSet
	case AnalysisStateInTriage.String():
		*a = AnalysisStateInTriage
	case AnalysisStateExploitable.String():
		*a = AnalysisStateExploitable
	case AnalysisStateNotAffected.String():
		*a = AnalysisStateNotAffected
	case AnalysisStateFalsePositive.String():
		*a = AnalysisStateFalsePositive
	default:
		return fmt.Errorf("invalid value: %s", val)
	}

	return nil
}

func (a AnalysisState) String() string {
	switch a {
	case AnalysisStateNotSet:
		return "NOT_SET"
	case AnalysisStateInTriage:
		return "IN_TRIAGE"
	case AnalysisStateExploitable:
		return "EXPLOITABLE"
	case AnalysisStateNotAffected:
		return "NOT_AFFECTED"
	case AnalysisStateFalsePositive:
		return "FALSE_POSITIVE"
	default:
		panic(a)
	}
}

type Analysis struct {
	Comments   []AnalysisComment `json:"comments"`
	State      AnalysisState     `json:"analysisState"`
	Suppressed bool              `json:"isSuppressed"`
}

type AnalysisComment struct {
	Comment   string `json:"comment"`
	Commenter string `json:"commenter"`
	Timestamp string `json:"timestamp"`
}

type AnalysisRequest struct {
	Component     uuid.UUID     `json:"component"`
	Project       uuid.UUID     `json:"project"`
	Vulnerability uuid.UUID     `json:"vulnerability"`
	Comment       string        `json:"comment,omitempty"`
	State         AnalysisState `json:"analysisState,omitempty"`
	Suppressed    bool          `json:"isSuppressed"`
}

func (c Client) GetAnalysis(ctx context.Context, component, project, vulnerability uuid.UUID) (*Analysis, error) {
	params := map[string]string{
		"component":     component.String(),
		"project":       project.String(),
		"vulnerability": vulnerability.String(),
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/analysis", withParams(params))
	if err != nil {
		return nil, err
	}

	var analysis Analysis
	_, err = c.doRequest(req, &analysis)
	if err != nil {
		return nil, err
	}

	return &analysis, nil
}

func (c Client) CreateAnalysis(ctx context.Context, analysisReq AnalysisRequest) (*Analysis, error) {
	req, err := c.newRequest(ctx, http.MethodPut, "/api/v1/analysis", withBody(analysisReq))
	if err != nil {
		return nil, err
	}

	var analysis Analysis
	_, err = c.doRequest(req, &analysis)
	if err != nil {
		return nil, err
	}

	return &analysis, nil
}
