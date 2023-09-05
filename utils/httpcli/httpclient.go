package httpcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/cwloo/gonet/logs"
)

func New(timeout int) *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{
		Jar:     jar,
		Timeout: time.Duration(timeout) * time.Second,
		// Transport: &http.Transport{
		// 	DisableKeepAlives:     false,
		// 	TLSHandshakeTimeout:   5 * time.Second,
		// 	IdleConnTimeout:       5 * time.Second,
		// 	ResponseHeaderTimeout: 5 * time.Second,
		// 	ExpectContinueTimeout: 5 * time.Second,
		// },
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after %v redirects", len(via))
			}
			return nil
		},
	}
}

func Get(url string, timeout int, c ...*http.Client) ([]byte, error) {
	logs.Infof("%v", url)
	switch len(c) > 0 {
	case true:
		return get(url, c[0])
	default:
		return get(url, &http.Client{Timeout: time.Duration(timeout) * time.Second})
	}
}

// application/json; charset=utf-8
func Post(url string, data any, timeout int, c ...*http.Client) ([]byte, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	logs.Infof("%v %v", url, string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("content-type", "application/json; charset=utf-8")
	switch len(c) > 0 {
	case true:
		return do(req, c[0])
	default:
		return do(req, &http.Client{Timeout: time.Duration(timeout) * time.Second})
	}
}

func PostReturn(url string, input, output any, timeout int, c ...*http.Client) error {
	b, err := Post(url, input, timeout, c...)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, output); err != nil {
		return err
	}
	return nil
}

func Check(url string, timeout int) bool {
	logs.Infof("%v %v", "GET", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", "bytes=0-1")

	rsp, err := (&http.Client{Timeout: time.Duration(timeout) * time.Second}).Do(req)
	if err != nil {
		return false
	}
	defer rsp.Body.Close()
	logs.Infof("%s", rsp.StatusCode)
	if rsp.StatusCode >= 200 && rsp.StatusCode < 300 {
		return true
	}
	return false
}

func get(url string, c *http.Client) ([]byte, error) {
	rsp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func do(req *http.Request, c *http.Client) ([]byte, error) {
	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
