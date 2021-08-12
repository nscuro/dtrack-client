package dtrack

type ProjectProperty struct {
	Group       string `json:"groupName"`
	Name        string `json:"propertyName"`
	Value       string `json:"propertyValue"`
	Type        string `json:"propertyType"`
	Description string `json:"description"`
}
