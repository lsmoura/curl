package curl

import (
	"context"
	"io"
)

type Request struct {
	ctx    context.Context
	url    string
	method string
	body   io.Reader

	FollowRedirects bool // Follow redirects
	Verbose         bool // Make the operation more talkative
	Out             io.Writer
}

// NewRequest creates a new http request with the given method and url
func NewRequest(method string, url string, body io.Reader) *Request {
	return NewRequestWithContext(context.Background(), method, url, body)
}

// NewRequestWithContext creates a new http request with the given method, url and context
func NewRequestWithContext(ctx context.Context, method string, url string, body io.Reader) *Request {
	return &Request{
		ctx:    ctx,
		method: method,
		url:    url,
		body:   body,
	}
}
