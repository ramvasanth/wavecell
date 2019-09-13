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
	SingleMessagePath = "sms/v1/"
)

// HTTPInterface helps wavecell tests
type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client manages requests to wavecell
type Client struct {
	BaseURL    string
	AuthKey    string
	ClientID   string
	HTTPClient HTTPInterface
}

func ClientWithAuthKey(key, clientID string) *Client {
	return &Client{
		BaseURL:    "https://api.wavecell.com/",
		ClientID:   clientID,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
		AuthKey:    key,
	}
}

// SingleMessage sends one message to one recipient
func (c Client) SingleMessage(m Message) (r Response, err error) {
	if err = m.Validate(); err != nil {
		return
	}
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	r, err = c.defaultRequest(b, SingleMessagePath+c.ClientID+"/single")
	return
}

func (c Client) defaultRequest(b []byte, path string) (r Response, err error) {
	req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewBuffer(b))
	if err != nil {
		return
	}
	if c.AuthKey != "" {
		req.Header.Add("Authorization", "Bearer "+c.AuthKey)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
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
