package notification

import "github.com/google/uuid"

type Component struct {
	UUID    uuid.UUID `json:"uuid"`
	Group   string    `json:"group"`
	Name    string    `json:"name"`
	Version string    `json:"version"`
	MD5     string    `json:"md5"`
	SHA1    string    `json:"sha1"`
	SHA256  string    `json:"sha256"`
	SHA512  string    `json:"sha512"`
	PURL    string    `json:"purl"`
}

type Policy struct {
	UUID           uuid.UUID `json:"uuid"`
	Name           string    `json:"name"`
	ViolationState string    `json:"violationState"`
}

type PolicyCondition struct {
	UUID     uuid.UUID `json:"uuid"`
	Subject  string    `json:"subject"`
	Operator string    `json:"operator"`
	Value    string    `json:"value"`
	Policy   Policy    `json:"policy"`
}

type PolicyViolation struct {
	UUID            uuid.UUID       `json:"uuid"`
	Type            string          `json:"type"`
	Timestamp       string          `json:"timestamp"`
	PolicyCondition PolicyCondition `json:"policyCondition"`
}

type Project struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	PURL        string    `json:"purl"`
	Tags        string    `json:"tags"`
}

type Vulnerability struct {
	UUID           uuid.UUID `json:"uuid"`
	VulnID         string    `json:"vulnId"`
	Source         string    `json:"source"`
	Title          string    `json:"title"`
	SubTitle       string    `json:"subtitle"`
	Description    string    `json:"description"`
	Recommendation string    `json:"recommendation"`
	CVSSV2         float32   `json:"cvssv2"`
	CVSSV3         float32   `json:"cvssv3"`
	Severity       string    `json:"severity"`
}
