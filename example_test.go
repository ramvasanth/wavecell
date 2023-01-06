package wavecell

import (
	"context"
	"log"
	"net/http"
	"time"
)

func ExampleNew() {
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
