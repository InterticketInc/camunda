package camunda

// ResponseCount a count response
type ResponseCount struct {
	Count int `json:"count"`
}

// ResLink a link struct
type ResLink struct {
	Method string `json:"method"`
	Href   string `json:"href"`
	Rel    string `json:"rel"`
}
