# wavecell

Wavecell API client library in Go

## Usage

To initiate a client, you should use the `wavecell.ClientWithAuthKey` func. This will returns a pointer to `wavecell.Client` and allows to you use features from service.

### Sending a single message

The func needs a `wavecell.Message` struct. That struct consists of the following attributes:

| Attribute | Type | Description |
|-----------|------|-------------|
| From | string | Represents a sender ID which can be alphanumeric or numeric |
| To | string | Message destination address |
| Text | string | Text of the message that will be sent |

It has a func to validate the `From` and `To` attributes, according to wavecell docs, and it is used into all funcs that make a request to the service. The following code is a basic example of the validate func:

```go
package main

import (
    "log"

    "github.com/ramvasanth/wavecell"
)

func main() {
    m := wavecell.Message{
        From: "Company", // or company number
        To:   "41793026727",
        Text: "This is an example of the body of the SMS message",
    }
    err := m.Validate()
    if err != nil {
        log.Fatalf("wavecell message error: %v", err)
    }
}
```

Finally, the following code is a full example to send a single message to wavecell service:

```go
package main

import (
    "fmt"
    "log"

    "github.com/ramvasanth/wavecell"
)

func main() {
    client := wavecell.ClientWithBasicAuth("foo", "bar")
    r, err := client.SingleMessage(m) // "m" refers to the variable from the previous example
    if err != nil {
        log.Fatalf("wavecell error: %v", err)
    }
    fmt.Printf("wavecell response: %v", r)
}
```