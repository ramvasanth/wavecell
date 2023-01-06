package wavecell

import (
	"github.com/shopspring/decimal"
)

// CallbackData specifies the callback data submitted by the vendor.
type CallbackData struct {
	Namespace   string              `json:"namespace"`
	EventType   string              `json:"eventType"`
	Description string              `json:"description"`
	Payload     CallbackDataPayload `json:"payload"`
}

// CallbackDataPayload specifies the message payload of callback data submitted by the vendor.
type CallbackDataPayload struct {
	UMID            string        `json:"umid"`
	BatchID         string        `json:"batchId"`
	ClientMessageID string        `json:"clientMessageId"`
	ClientBatchID   string        `json:"clientBatchId"`
	SubAccountID    string        `json:"subAccountId"`
	Source          string        `json:"source"`
	Destination     string        `json:"destination"`
	Status          MessageStatus `json:"status"`
	Price           MessagePrice  `json:"price"`
	SmsCount        int           `json:"smsCount"`
}

// MessageStatus specifies the message status.
type MessageStatus struct {
	State        string `json:"state"`
	Detail       string `json:"detail"`
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// MessagePrice specifies the message price.
type MessagePrice struct {
	Total    *decimal.Decimal `json:"total"`
	PerSMS   decimal.Decimal  `json:"perSms"`
	Currency *string          `json:"currency"`
}

// List of Status Callback
const (
	StatusCallbackUnknown     = `unknown`
	StatusCallbackQueued      = `queued`
	StatusCallbackFailed      = `failed`
	StatusCallbackSent        = `sent`
	StatusCallbackDelivered   = `delivered`
	StatusCallbackUndelivered = `undelivered`
	StatusCallbackRead        = `read`
	StatusCallbackOk          = `ok`
	StatusCallbackError       = `error`
	StatusCallbackRejected    = `rejected`
)

// List of Detail Callback
const (
	DetailCallbackDeliveredToOperator    = `delivered_to_operator`
	DetailCallbackDeliveredToRecipient   = `delivered_to_recipient`
	DetailCallbackRejectedByOperator     = `rejected_by_operator`
	DetailCallbackUndeliveredToRecipient = `undelivered_to_recipient`
)
