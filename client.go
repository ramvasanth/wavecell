package wavecell

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	//SingleMessagePath for sending a single message
	SINGLE_SMS_ENDPOINT = "sms/v1/%s/single"
	CONNECTION_TIME_OUT = 15
)

// HTTPInterface helps wavecell tests
type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client manages requests to wavecell
type Config struct {
	BaseURL    string
	AuthKey    string
	ClientID   string
}

type Sender struct{
	Config Config
	HTTPClient HTTPInterface
}

func New(config Config) *Sender {
	return &Sender{
		Config: config,
		HTTPClient: &http.Client{Timeout: CONNECTION_TIME_OUT * time.Second},
	}
}

// SingleMessage sends one message to one recipient
func (s *Sender) SingleMessage(m Message) (r Response, err error) {
	if err = m.Validate(); err != nil {
		return
	}
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	path := fmt.Sprintf(SINGLE_SMS_ENDPOINT, s.Config.ClientID)
	r, err = s.defaultRequest(b, path)
	return
}

func (s *Sender) defaultRequest(b []byte, path string) (r Response, err error) {
	req, err := http.NewRequest(http.MethodPost, s.Config.BaseURL+path, bytes.NewBuffer(b))
	if err != nil {
		return
	}
	if s.Config.AuthKey != "" {
		req.Header.Add("Authorization", "Bearer "+s.Config.AuthKey)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		err = json.NewDecoder(resp.Body).Decode(&r)
		return
	}
	bEr, _ := ioutil.ReadAll(resp.Body)
	err = fmt.Errorf("received status code: %d, error: %s,for more information on this error refer to this link:https://developer.wavecell.com/v1/sms-api/api-send-sms", resp.StatusCode, bEr)
	return
}
