# wavecell

[![Go Reference](https://pkg.go.dev/badge/github.com/flip-id/wavecell.svg)](https://pkg.go.dev/github.com/flip-id/wavecell)
[![Go Report Card](https://goreportcard.com/badge/github.com/flip-id/wavecell)](https://goreportcard.com/report/github.com/flip-id/wavecell)

Wavecell API client library in Go.

# How to Use

```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	c, err := New(
		WithAPIKey("YOUR_API_KEY"),
		WithTimeout(1*time.Minute),
		WithSubAccountID("SUB_ACCOUNT_ID"),
		WithClient(http.DefaultClient),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.SendSMSV1(context.Background(), &RequestSendSMS{
		Destination: "+62101010101",
		Text:        "Hello!",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response: %+v", resp)
}

```