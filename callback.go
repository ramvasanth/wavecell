package wavecell

type CallbackData struct {
	BatchID         string `json:"batchId"`
	ClientBatchID   string `json:"clientBatchId"`
	ClientMessageID string `json:"clientMessageId"`
	Destination     string `json:"destination"`
	Error           string `json:"error"`
	ErrorCode       int64  `json:"errorCode"`
	Price           struct {
		Currency string  `json:"currency"`
		PerSms   float64 `json:"perSms"`
		Total    float64 `json:"total"`
	} `json:"price"`
	SmsCount     int64  `json:"smsCount"`
	Source       string `json:"source"`
	Status       string `json:"status"`
	StatusCode   int64  `json:"statusCode"`
	SubAccountID string `json:"subAccountId"`
	Timestamp    string `json:"timestamp"`
	Umid         string `json:"umid"`
	Version      int64  `json:"version"`
}
