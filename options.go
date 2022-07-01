package wavecell

import (
	"github.com/gojek/heimdall/v7/hystrix"
	"net/http"
	"strings"
	"time"

	"github.com/gojek/heimdall/v7"
)

// DefaultTimeout is the default timeout of the Wavecell API.
const DefaultTimeout = 30 * time.Second

// BaseURL is the base URL of the Wavecell API.
const BaseURL = "https://sms.8x8.com"

// Option is the option for the client.
type Option struct {
	BaseURL        string
	APIKey         string
	SubAccountID   string
	Client         heimdall.Doer
	Timeout        time.Duration
	HystrixOptions []hystrix.Option
	client         *hystrix.Client
}

func (o *Option) Assign(opts ...FnOption) *Option {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Clone returns a clone of the option.
func (o *Option) Clone() *Option {
	opt := *o
	return &opt
}

// Default is the default option for the client.
func (o *Option) Default() *Option {
	if o.BaseURL == "" {
		o.BaseURL = BaseURL
	}

	o.BaseURL = strings.TrimRight(o.BaseURL, "/")
	if o.Client == nil {
		o.Client = http.DefaultClient
	}

	if o.Timeout < DefaultTimeout {
		o.Timeout = DefaultTimeout
	}

	opts := append([]hystrix.Option{
		hystrix.WithHTTPTimeout(o.Timeout),
		hystrix.WithHystrixTimeout(o.Timeout),
		hystrix.WithHTTPClient(o.Client),
	},
		o.HystrixOptions...,
	)
	o.client = hystrix.NewClient(
		opts...,
	)
	return o
}

// Validate validates the option.
func (o *Option) Validate() (err error) {
	if o.APIKey == "" {
		err = ErrEmptyAPIKEY
		return
	}

	if o.SubAccountID == "" {
		err = ErrEmptySubAccountID
	}
	return
}

// FnOption is the functional option to set the Option.
type FnOption func(o *Option)

// WithBaseURL sets the base URL of the Wavecell API.
func WithBaseURL(u string) FnOption {
	return func(o *Option) {
		o.BaseURL = u
	}
}

// WithAPIKey sets the API key of the Wavecell API.
func WithAPIKey(key string) FnOption {
	return func(o *Option) {
		o.APIKey = key
	}
}

// WithSubAccountID sets the sub account ID of the Wavecell API.
func WithSubAccountID(id string) FnOption {
	return func(o *Option) {
		o.SubAccountID = id
	}
}

// WithClient sets the client of the Wavecell API.
func WithClient(c heimdall.Doer) FnOption {
	return func(o *Option) {
		o.Client = c
	}
}

// WithTimeout sets the timeout of the Wavecell API.
func WithTimeout(t time.Duration) FnOption {
	return func(o *Option) {
		o.Timeout = t
	}
}

// WithHystrixOptions sets the hystrix options of the Wavecell API.
func WithHystrixOptions(opts ...hystrix.Option) FnOption {
	return func(o *Option) {
		o.HystrixOptions = opts
	}
}
