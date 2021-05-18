package http

import (
	"context"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Caller struct {
	client *http.Client
}

func NewCaller() *Caller {
	return &Caller{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Caller) Get(ctx context.Context, url string, params url.Values, token string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/http.Caller.Get create request fail")
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Add("token", "token")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/http.Caller.Get do request fail")
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("pkg/http.Caller.Get response error : " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "pkg/http.Caller.Get read response error")
	}
	return body, nil
}
