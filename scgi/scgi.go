package scgi

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

type Client struct {
	Network string
	Address string
}

type Header struct {
	Name  string
	Value string
}

type Request struct {
	Headers []Header
	Body    []byte
}

type Response Request

func New(network string, address string) (c *Client) {
	c = &Client{Network: network, Address: address}
	return
}

func (c *Client) RoundTrip(req *http.Request) (*http.Response, error) {
	var conn net.Conn
	var err error
	resp := &http.Response{}

	switch c.Network {
	case "tcp", "unix":
		conn, err = net.Dial(c.Network, c.Address)
	}

	if err == nil {
		defer conn.Close()
		headers := [][2]string{
			{"CONTENT_LENGTH", strconv.FormatInt(req.ContentLength, 10)},
			{"SCGI", "1"},
			{"REQUEST_METHOD", req.Method},
			{"REQUEST_URI", req.RequestURI},
			{"SERVER_PROTOCOL", req.Proto},
		}

		var headersBytes bytes.Buffer
		for _, h := range headers {
			headersBytes.WriteString(h[0])
			headersBytes.WriteRune('\x00')
			headersBytes.WriteString(h[1])
			headersBytes.WriteRune('\x00')
		}

		conn.Write([]byte(strconv.Itoa(headersBytes.Len())))
		conn.Write([]byte(":"))
		conn.Write(headersBytes.Bytes())
		conn.Write([]byte(","))

		body, _ := ioutil.ReadAll(req.Body)
		conn.Write(body)

		defer req.Body.Close()

		// TODO: parse headers properly?
		reader := bufio.NewReader(conn)
		var line []byte
		for {
			line, _, _ = reader.ReadLine()
			if len(line) == 0 {
				break
			}
		}
		body, _ = ioutil.ReadAll(reader)

		resp.Status = http.StatusText(200)
		resp.StatusCode = 200
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		resp.Request = req
	}
	return resp, err
}
