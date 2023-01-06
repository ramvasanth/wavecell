package wavecell

import (
	"context"
	"fmt"
	"github.com/fairyhunter13/iso8601/v2"
	"github.com/gojek/heimdall/v7/hystrix"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []FnOption
	}
	tests := []struct {
		name          string
		args          func() args
		wantNilClient bool
		wantErr       bool
	}{
		{
			name: "API Key is Empty",
			args: func() args {
				return args{
					opts: []FnOption{
						WithAPIKey(""),
					},
				}
			},
			wantNilClient: true,
			wantErr:       true,
		},
		{
			name: "Empty Sub Account ID",
			args: func() args {
				return args{
					opts: []FnOption{
						WithSubAccountID(""),
						WithBaseURL("https://sms.8x8.com"),
						WithAPIKey("123456789"),
						WithSubAccountID(""),
						WithTimeout(time.Second * 30),
						WithClient(http.DefaultClient),
					},
				}
			},
			wantNilClient: true,
			wantErr:       true,
		},
		{
			name: "Empty Base URL",
			args: func() args {
				return args{
					opts: []FnOption{
						WithAPIKey("123456789"),
						WithSubAccountID("123456789"),
						WithTimeout(time.Second * 10),
						WithClient(http.DefaultClient),
					},
				}
			},
			wantNilClient: false,
			wantErr:       false,
		},
		{
			name: "Timeout Less Than 30 Seconds",
			args: func() args {
				return args{
					opts: []FnOption{
						WithBaseURL("https://sms.8x8.com"),
						WithAPIKey("123456789"),
						WithSubAccountID("123456789"),
						WithTimeout(time.Second * 10),
						WithClient(http.DefaultClient),
					},
				}
			},
			wantNilClient: false,
			wantErr:       false,
		},
		{
			name: "Nil HTTP Client",
			args: func() args {
				return args{
					opts: []FnOption{
						WithBaseURL("https://sms.8x8.com"),
						WithAPIKey("123456789"),
						WithSubAccountID("123456789"),
						WithTimeout(time.Second * 30),
						WithClient(nil),
					},
				}
			},
			wantNilClient: false,
			wantErr:       false,
		},
		{
			name: "Success Initializing Client",
			args: func() args {
				return args{
					opts: []FnOption{
						WithBaseURL("https://sms.8x8.com"),
						WithAPIKey("123456789"),
						WithSubAccountID("123456789"),
						WithTimeout(time.Second * 30),
						WithClient(http.DefaultClient),
						WithHystrixOptions(
							hystrix.WithHystrixTimeout(time.Second * 30),
						),
					},
				}
			},
			wantNilClient: false,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args()
			gotC, err := New(args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantNilClient {
				assert.Nil(t, gotC)
			}
		})
	}
}

func Test_client_SendSMSV1(t *testing.T) {
	timeCollections := []iso8601.Time{
		{time.Now().UTC()},
	}
	type fields struct {
		opt *Option
	}
	type args struct {
		req *RequestSendSMS
	}
	tests := []struct {
		name      string
		fields    func() fields
		args      func() args
		wantResp  func() *ResponseSendSMS
		wantErr   bool
		responder func()
	}{
		{
			name: "Invalid URL",
			fields: func() fields {
				return fields{
					opt: (new(Option)).Assign(
						WithBaseURL("://sms.8x8.com"),
						WithSubAccountID("test"),
						WithAPIKey("test"),
					).Default(),
				}
			},
			args: func() args {
				req := RequestSendSMS{
					Destination: "+6281212121212",
					Text:        "Test Message",
				}
				return args{req: &req}
			},
			wantResp: func() *ResponseSendSMS {
				return nil
			},
			wantErr: true,
		},
		{
			name: "Server Unreached",
			fields: func() fields {
				return fields{
					opt: (new(Option)).Assign(
						WithBaseURL("http://test.local"),
						WithSubAccountID("test"),
						WithAPIKey("test"),
					).Default(),
				}
			},
			args: func() args {
				req := RequestSendSMS{
					Destination: "+6281212121212",
					Text:        "Test Message",
				}
				return args{req: &req}
			},
			wantResp: func() *ResponseSendSMS {
				return nil
			},
			wantErr: true,
			responder: func() {
			},
		},
		{
			name: "Error Response from Server - Invalid Number",
			fields: func() fields {
				return fields{
					opt: (new(Option)).Assign(
						WithBaseURL("http://sms.8x8.com"),
						WithSubAccountID("test"),
						WithAPIKey("test"),
					).Default(),
				}
			},
			args: func() args {
				req := RequestSendSMS{
					Destination: "+62851xxx1121",
					Text:        "Test Message",
				}
				return args{req: &req}
			},
			wantResp: func() *ResponseSendSMS {
				return nil
			},
			wantErr: true,
			responder: func() {
				httpmock.RegisterResponder(
					http.MethodPost,
					"http://sms.8x8.com"+fmt.Sprintf(URLSendSMS, "test"),
					httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, &ResponseError{
						Code:      1002,
						Message:   "Invalid MSISDN format (not E.164 international number)",
						ErrorID:   "b4478860-b76c-e811-814e-022a35cc1c71",
						Timestamp: timeCollections[0],
					}),
				)
			},
		},
		{
			name: "Success Response from Server",
			fields: func() fields {
				return fields{
					opt: (new(Option)).Assign(
						WithBaseURL("http://sms.8x8.com"),
						WithSubAccountID("test"),
						WithAPIKey("test"),
					).Default(),
				}
			},
			args: func() args {
				req := RequestSendSMS{
					Destination: "+628511011121",
					Text:        "Test Message",
				}
				return args{req: &req}
			},
			wantResp: func() *ResponseSendSMS {
				return &ResponseSendSMS{
					UmID:            "bda3d56d-1424-e711-813c-06ed3428fe67",
					ClientMessageID: "client-message-id",
					Destination:     "+628511011121",
					Encoding:        "GSM7",
					Status: ResponseSendSMSStatus{
						Code:        "QUEUED",
						Description: "SMS is accepted and queued for processing",
					},
				}
			},
			wantErr: false,
			responder: func() {
				httpmock.RegisterResponder(
					http.MethodPost,
					"http://sms.8x8.com"+fmt.Sprintf(URLSendSMS, "test"),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, &ResponseSendSMS{
						UmID:            "bda3d56d-1424-e711-813c-06ed3428fe67",
						ClientMessageID: "client-message-id",
						Destination:     "+628511011121",
						Encoding:        "GSM7",
						Status: ResponseSendSMSStatus{
							Code:        "QUEUED",
							Description: "SMS is accepted and queued for processing",
						},
					}),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			if tt.responder != nil {
				tt.responder()
			}

			fields := tt.fields()
			args := tt.args()
			wantResp := tt.wantResp()
			s := &client{
				opt: fields.opt,
			}
			gotResp, err := s.SendSMSV1(context.TODO(), args.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			}

			assert.EqualValuesf(t, wantResp, gotResp, "SendSMSV1(%v)", args.req)
		})
	}
}
