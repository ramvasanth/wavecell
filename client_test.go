package wavecell

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type HTTPClientMock struct {
	Response *http.Response
}

func (c *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	c.Response.StatusCode = 200
	return c.Response, nil
}

func TestForResponseMessages(t *testing.T) {
	tests := []struct {
		reference string
		body      string
		data      Response
	}{
		{
			reference: "#1",
			body:      `{}`,
			data:      Response{},
		},
		{
			reference: "#2",
			body:      `{"messages": []}`,
			data:      Response{},
		},
		{
			reference: "#3",
			body:      `{"destination": "41793026727"}`,
			data: Response{
				Destination: "41793026727",
			},
		},
		{
			reference: "#4",
			body:      `{"destination": "41793026727"}`,
			data: Response{
				Destination: "41793026727",
			},
		},
		{
			reference: "#5",
			body: `{"umid": "bda3d56d-1424-e711-813c-06ed3428fe67", "clientMessageId": "1234", "destination": "41793026727", "encoding": "GSM7",
			"status": {
			"code": "QUEUED",
			"description": "SMS is accepted and queued for processing"
			} }`,
			data: Response{
				Destination:     "41793026727",
				ClientMessageID: "1234",
				UMID:            "bda3d56d-1424-e711-813c-06ed3428fe67",
				Encoding:        "GSM7",
				Status: struct {
					Code        string `json:"code"`
					Description string `json:"description"`
				}{"QUEUED", "SMS is accepted and queued for processing"},
			},
		},
	}

	var message = Message{
		From: "company",
		To:   "442071838750",
		Text: "Foo bar",
	}

	client := ClientWithAuthKey("foo", "bar")
	for _, test := range tests {
		test := test
		t.Run(test.reference, func(t *testing.T) {
			client.HTTPClient = &HTTPClientMock{
				Response: &http.Response{
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.body)),
					StatusCode: 200,
				},
			}
			r, err := client.SingleMessage(message)
			if err != nil {
				t.Errorf("Error: unexpected error was returned (%s)", err)
			}
			if !reflect.DeepEqual(r, test.data) {
				t.Errorf("expected '%v', got '%v'", test.data, r)
			}
		})
	}
}

func TestForSingleMessageError(t *testing.T) {
	tests := []struct {
		reference string
		message   Message
		err       error
	}{
		{
			reference: "#1",
			message:   Message{},
			err:       ErrForFromNonAlphanumeric,
		},
	}

	client := ClientWithAuthKey("key", "cliend_id")
	client.HTTPClient = &HTTPClientMock{
		Response: &http.Response{
			Body: ioutil.NopCloser(nil),
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.reference, func(t *testing.T) {
			if _, err := client.SingleMessage(test.message); err != test.err {
				t.Errorf("expected '%v', got '%v'", test.err, err)
			}
		})
	}
}
