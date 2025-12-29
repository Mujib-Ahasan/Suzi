package attacks

import "time"

type Options struct {
	URL         string
	Rate        int
	Requests    int
	Timeout     time.Duration
	Type        string
	Method      string
	Body        []byte
	ContentType string
}
