package wavecell

// ResponseStatus ...
type Response struct {
	UMID            string `json:"umid"`
	ClientMessageID string `json:"clientMessageId,omitempty"`
	Destination     string `json:"destination"`
	Encoding        string `json:"encoding"`
	Status          struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"status"`
}
