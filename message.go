package wavecell

import "regexp"

// Message contains the body request
type Message struct {
	From            string `json:"source,omitempty"`
	To              string `json:"destination"`
	Text            string `json:"text"`
	ClientMessageID string `json:"clientMessageId"`
	Encoding        string `json:"encoding"`
	DlrCallbackUrl  string `json:"dlrCallbackUrl"`
}

// Validate validates the body request values
func (m Message) Validate() (err error) {
	if err = m.validateFromValue(); err != nil {
		return
	}

	err = m.validateToValue()
	return
}

func (m Message) validateFromValue() (err error) {
	if isNumeric(m.From) && !isValidRange(m.From, 3, 14) {
		err = ErrForFromNonAlphanumeric
		return
	}
	if !isValidRange(m.From, 3, 13) {
		err = ErrForFromAlphanumeric
		return
	}
	return
}

func (m Message) validateToValue() (err error) {
	if m.To == "" {
		return
	}
	if isNumeric(m.To) && !isValidRange(m.To, 3, 14) {
		err = ErrForToNonAlphanumeric
		return
	}
	return
}

func isNumeric(s string) bool {
	return regexp.MustCompile(`^[\d]*$`).MatchString(s)
}

func isValidRange(s string, a, b int) bool {
	l := len(s)
	return l > a && l <= b
}
