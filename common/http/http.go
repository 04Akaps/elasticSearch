package http

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
	"time"
)

const (
	_paramLength = "get param length and request length"
)

type Client struct {
	client *resty.Client
}

func NewClient(timeOut time.Duration, baseURL string) *Client {

	client := resty.New().
		SetJSONMarshaler(sonic.Marshal).
		SetJSONUnmarshaler(sonic.Unmarshal).
		SetTimeout(timeOut).
		SetBaseURL(baseURL)

	return &Client{client}
}

func (c *Client) SetHeader(headers map[string]string) {
	c.client.SetHeaders(headers)
}

func (c *Client) SetAuth(schema string, v string) {
	// default schema is bearer
	if schema != "" {
		c.client.SetAuthScheme(schema)
	}
	c.client.SetAuthToken(v)
}

func (c *Client) POST(url string, req, resp interface{}) error {
	body, err := c.client.JSONMarshal(req)

	if err != nil {
		return err
	}

	_, err = c.client.R().
		SetBody(body).
		SetResult(&resp).
		Post(url)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GET(url string, paramName, req []string, resp interface{}) error {
	if len(paramName) != len(req) {
		return errors.New(_paramLength)
	}

	queryParam := make(map[string]string, len(paramName))

	for i, v := range paramName {
		queryParam[v] = req[i]
	}

	_, err := c.client.R().
		SetQueryParams(queryParam).
		SetResult(&resp).
		Get(url)

	if err != nil {
		return err
	}

	return nil
}
