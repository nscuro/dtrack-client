package dtrack

import (
	"encoding/json"
	"fmt"
)

type Notification struct {
	Level     string
	Scope     string
	Group     string
	Timestamp string
	Title     string
	Content   string
	Subject   interface{}
}

type notificationJSON struct {
	Level     string          `json:"level"`
	Scope     string          `json:"scope"`
	Group     string          `json:"group"`
	Timestamp string          `json:"timestamp"`
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	Subject   json.RawMessage `json:"subject"`
}

func ParseNotification(notification []byte) (*Notification, error) {
	nw := struct {
		Notification notificationJSON `json:"notification"`
	}{}
	err := json.Unmarshal(notification, &nw)
	if err != nil {
		return nil, err
	}

	var subject interface{}
	switch nw.Notification.Group {
	case "BOM_CONSUMED":
		fallthrough
	case "BOM_PROCESSED":
		subject = &BOMSubject{}
	case "NEW_VULNERABLE_DEPENDENCY":
		subject = &NewVulnerableDependencySubject{}
	case "NEW_VULNERABILITY":
		subject = &NewVulnerabilitySubject{}
	default:
		return nil, fmt.Errorf("unknown notification group %s", nw.Notification.Group)
	}

	err = json.Unmarshal(nw.Notification.Subject, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal subject: %w", err)
	}

	return &Notification{
		Level:     nw.Notification.Level,
		Scope:     nw.Notification.Scope,
		Group:     nw.Notification.Group,
		Timestamp: nw.Notification.Timestamp,
		Title:     nw.Notification.Title,
		Content:   nw.Notification.Content,
		Subject:   subject,
	}, nil
}

type BOMSubject struct {
	BOM struct {
		Content     string `json:"content"`
		Format      string `json:"format"`
		SpecVersion string `json:"specVersion"`
	} `json:"bom"`
	Project Project `json:"project"`
}

type NewVulnerableDependencySubject struct {
	Component       Component       `json:"component"`
	Project         Project         `json:"project"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}

type NewVulnerabilitySubject struct {
	AffectedProjects []Project     `json:"affectedProjects"`
	Component        Component     `json:"component"`
	Vulnerability    Vulnerability `json:"vulnerability"`
}
