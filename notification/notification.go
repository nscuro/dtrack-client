package notification

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const (
	GroupBOMConsumed             = "BOM_CONSUMED"
	GroupBOMProcessed            = "BOM_PROCESSED"
	GroupNewVulnerableDependency = "NEW_VULNERABLE_DEPENDENCY"
	GroupNewVulnerability        = "NEW_VULNERABILITY"
	GroupPolicyViolation         = "POLICY_VIOLATION"
	GroupVEXConsumed             = "VEX_CONSUMED"
	GroupVEXProcessed            = "VEX_PROCESSED"

	LevelError         = "ERROR"
	LevelInformational = "INFORMATIONAL"
	LevelWarning       = "WARNING"

	ScopeSystem    = "SYSTEM"
	ScopePortfolio = "PORTFOLIO"
)

type Notification struct {
	Level     string
	Scope     string
	Group     string
	Timestamp time.Time
	Title     string
	Content   string
	Subject   interface{}
}

type notificationJSON struct {
	Level     string          `json:"level"`
	Scope     string          `json:"scope"`
	Group     string          `json:"group"`
	Timestamp timestampJSON   `json:"timestamp"`
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	Subject   json.RawMessage `json:"subject"`
}

type timestampJSON struct {
	time.Time
}

func (t *timestampJSON) UnmarshalJSON(data []byte) (err error) {
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil || str == "" {
		return
	}

	parsedTime, err := time.Parse(`2006-01-02T15:04:05.99`, str)
	if err != nil {
		return
	}

	*t = timestampJSON{parsedTime}
	return
}

type notificationWrapperJSON struct {
	Notification notificationJSON `json:"notification"`
}

// Parse parses a notification.
func Parse(reader io.Reader) (n Notification, err error) {
	wrapper := notificationWrapperJSON{}
	err = json.NewDecoder(reader).Decode(&wrapper)
	if err != nil {
		return
	}

	var subject interface{}
	switch wrapper.Notification.Group {
	case GroupBOMConsumed:
		fallthrough
	case GroupBOMProcessed:
		subject = &BOMSubject{}
	case GroupNewVulnerableDependency:
		subject = &NewVulnerableDependencySubject{}
	case GroupNewVulnerability:
		subject = &NewVulnerabilitySubject{}
	case GroupPolicyViolation:
		subject = &PolicyViolationSubject{}
	case GroupVEXConsumed:
		fallthrough
	case GroupVEXProcessed:
		subject = &VEXSubject{}
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
		Timestamp: wrapper.Notification.Timestamp.Time,
		Title:     wrapper.Notification.Title,
		Content:   wrapper.Notification.Content,
		Subject:   subject,
	}, nil
}
