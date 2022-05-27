package notification

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

type PolicyViolationSubject struct {
	Component       Component       `json:"component"`
	PolicyViolation PolicyViolation `json:"policyViolation"`
	Project         Project         `json:"project"`
}

type VEXSubject struct {
	VEX struct {
		Content     string `json:"content"`
		Format      string `json:"format"`
		SpecVersion string `json:"specVersion"`
	}
	Project Project `json:"project"`
}
