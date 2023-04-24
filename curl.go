package curl

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

var zeroDialer net.Dialer

type loggingConn struct {
	net.Conn

	writeBreak  bool
	afterHeader bool
}

func newLoggingConn(ctx context.Context, network string, address string) (net.Conn, error) {
	c, err := zeroDialer.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}

	return loggingConn{Conn: c}, nil
}

func (l loggingConn) Write(b []byte) (int, error) {
	if l.writeBreak {
		fmt.Fprintf(os.Stderr, "\n")
		l.writeBreak = false
	}
	data := string(b)
	for len(data) > 0 {
		idx := strings.Index(data, "\r\n")
		if idx == -1 {
			break
		}
		fmt.Fprintf(os.Stderr, "> %s\n", data[:idx])
		data = data[idx+2:]
	}

	if len(data) > 0 {
		fmt.Fprintf(os.Stderr, "> %s", data)
		l.writeBreak = true
	}

	return l.Conn.Write(b)
}

func (l loggingConn) Read(b []byte) (int, error) {
	n, err := l.Conn.Read(b)
	if err != nil {
		return n, err
	}

	if !l.afterHeader {
		data := string(b)

		for len(data) > 0 {
			idx := strings.Index(data, "\r\n")
			if idx == -1 {
				break
			}

			fmt.Fprintf(os.Stderr, "< %s\n", data[:idx])
			data = data[idx+2:]

			if idx == 0 {
				l.afterHeader = false
				break
			}
		}
	}
	return n, err
}

func clientForRequest(req *Request) *http.Client {
	client := http.Client{}

	if !req.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}

	}

	if req.Verbose {
		client.Transport = &http.Transport{
			DialContext: newLoggingConn,
		}
	}

	return &client
}

func Do(req *Request) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(req.ctx, req.method, req.url, req.body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	client := clientForRequest(req)

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close()

	out := req.Out
	if out == nil {
		out = os.Stdout
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	return data, nil
}
