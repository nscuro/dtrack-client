package dtrack

import (
	"github.com/google/uuid"
)

type Finding struct {
	Component     Component           `json:"component"`
	Vulnerability Vulnerability       `json:"vulnerability"`
	Attribution   *FindingAttribution `json:"attribution"`
	Analysis      *Analysis           `json:"analysis"`
	Matrix        string              `json:"matrix"`
}

type FindingAttribution struct {
	UUID                uuid.UUID `json:"uuid"`
	AnalyzerIdentity    string    `json:"analyzerIdentity"`
	AlternateIdentifier string    `json:"alternateIdentifier"`
	AttributedOn        string    `json:"attributedOn"`
	ReferenceURL        string    `json:"referenceUrl"`
}
