package wavecell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fairyhunter13/iso8601/v2"
	"github.com/fairyhunter13/phone"
	"github.com/fairyhunter13/pool"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// Client is the contract for the Wavecell client.
type Client interface {
	// SendSMSV1 sends one message to one recipient.
	// The resp here can be either *ResponseError, *ResponseSendSMS, or nil.
	// This method is based on the documentation at: https://developer.8x8.com/connect/reference/send-sms-single.
	SendSMSV1(ctx context.Context, req *RequestSendSMS) (resp *ResponseSendSMS, err error)
}

type client struct {
	opt *Option
}

// Assign assigns the opt to the client.
func (c *client) Assign(opt *Option) *client {
	if opt == nil {
		return c
	}

	c.opt = opt.Clone()
	return c
}

// New returns a new Sender struct.
func New(opts ...FnOption) (c Client, err error) {
	o := new(Option).Assign(opts...).Default()
	err = o.Validate()
	if err != nil {
		return
	}

	c = (new(client)).Assign(o)
	return
}

// ResponseSendSMS is the response struct for SendSMSV1.
type ResponseSendSMS struct {
	UmID            string                `json:"umid"`
	Destination     string                `json:"destination"`
	Status          ResponseSendSMSStatus `json:"status"`
	Encoding        string                `json:"encoding"`
	ClientMessageID string                `json:"clientMessageId,omitempty"`
}

// ResponseSendSMSStatus is the response struct for SendSMSV1.
type ResponseSendSMSStatus struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// RequestSendSMS is the request struct for SendSMSV1.
type RequestSendSMS struct {
	Destination     string        `json:"destination,omitempty" validate:"required"`
	Country         string        `json:"country,omitempty"`
	Source          string        `json:"source,omitempty"`
	ClientMessageID string        `json:"clientMessageId,omitempty"`
	Text            string        `json:"text,omitempty" validate:"required"`
	Encoding        string        `json:"encoding,omitempty"`
	Scheduled       *iso8601.Time `json:"scheduled,omitempty"`
	Expiry          *iso8601.Time `json:"expiry,omitempty"`
	DlrCallbackURL  string        `json:"dlrCallbackUrl,omitempty"`
	ClientIP        string        `json:"clientIp,omitempty"`
	Track           string        `json:"track,omitempty"`
}

// Normalize normalizes the request.
func (r *RequestSendSMS) Normalize() *RequestSendSMS {
	r.Destination = phone.NormalizeID(r.Destination, 0)
	return r
}

func (s *client) getReqBuffer(req interface{}) (buf *bytes.Buffer, err error) {
	buf = pool.GetBuffer()
	err = json.NewEncoder(buf).Encode(req)
	return
}

func (s *client) clean(buff *bytes.Buffer) {
	pool.Put(buff)
}

func (s *client) getFullURL(action string) string {
	return s.opt.BaseURL + action
}

// SendSMSV1 sends one message to one recipient.
// The resp here can be either *ResponseError, *ResponseSendSMS, or nil.
func (s *client) SendSMSV1(ctx context.Context, req *RequestSendSMS) (resp *ResponseSendSMS, err error) {
	buff, err := s.getReqBuffer(req.Normalize())
	if err != nil {
		return
	}
	defer s.clean(buff)

	endpoint := fmt.Sprintf(URLSendSMS, s.opt.SubAccountID)
	r, err := http.NewRequest(http.MethodPost, s.getFullURL(endpoint), buff)
	if err != nil {
		return
	}

	res, err := s.doRequest(ctx, r)
	if err != nil {
		return
	}
	defer func() {
		if res.Body != nil {
			_ = res.Body.Close()
		}
	}()

	decoder := json.NewDecoder(res.Body)
	if s.isResponseError(res) {
		var respErr ResponseError
		err = decoder.Decode(&respErr)
		if err != nil {
			return
		}

		err = &respErr
		return
	}

	err = decoder.Decode(&resp)
	return
}

func (s *client) isResponseError(resp *http.Response) bool {
	return resp.StatusCode >= http.StatusBadRequest
}

func (s *client) doRequest(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	if s.opt.APIKey != "" {
		req.Header.Add(fiber.HeaderAuthorization, HeaderAuthBearerValue+s.opt.APIKey)
	}

	req.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.Header.Add(fiber.HeaderAccept, fiber.MIMEApplicationJSON)
	req = req.WithContext(ctx)
	resp, err = s.opt.client.Do(req)
	return
}
