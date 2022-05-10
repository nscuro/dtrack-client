package dtrack

import (
	"encoding/json"
	"fmt"
	"io"
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

func ParseNotification(reader io.Reader) (n Notification, err error) {
	wrapper := struct {
		Notification notificationJSON `json:"notification"`
	}{}
	err = json.NewDecoder(reader).Decode(&wrapper)
	if err != nil {
		return
	}

	var subject interface{}
	switch wrapper.Notification.Group {
	case "BOM_CONSUMED":
		fallthrough
	case "BOM_PROCESSED":
		subject = &BOMSubject{}
	case "NEW_VULNERABLE_DEPENDENCY":
		subject = &NewVulnerableDependencySubject{}
	case "NEW_VULNERABILITY":
		subject = &NewVulnerabilitySubject{}
	default:
		err = fmt.Errorf("unknown notification group %s", wrapper.Notification.Group)
		return
	}

	err = json.Unmarshal(wrapper.Notification.Subject, subject)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal subject: %w", err)
		return
	}

	return Notification{
		Level:     wrapper.Notification.Level,
		Scope:     wrapper.Notification.Scope,
		Group:     wrapper.Notification.Group,
		Timestamp: wrapper.Notification.Timestamp,
		Title:     wrapper.Notification.Title,
		Content:   wrapper.Notification.Content,
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
