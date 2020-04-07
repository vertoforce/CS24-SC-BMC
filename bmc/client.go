package bmc

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Client to communicate with BCM
type Client struct {
	ip              string
	port            uint16
	phpsessidCookie *http.Cookie
	httpClient      *http.Client
	Certificates    []*x509.Certificate
	CipherSuite     uint16
}

// New create a new client, baseURL is for example `https://10.0.0.130`
func New(ctx context.Context, ip string, port uint16, username, password string) (*Client, error) {
	// Get client
	c := &Client{ip: ip, port: port}
	var err error
	c.httpClient, c.Certificates, c.CipherSuite, err = createHTTPClient(ctx, ip, port)
	if err != nil {
		return nil, err
	}

	// Do login request
	form := url.Values{}
	form.Add("quser", username)
	form.Add("qpass", password)
	req, err := c.buildRequest(ctx, "POST", fmt.Sprintf("https://%s:%d/cgi_bin/login.cgi", ip, port), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	if resp.Request.URL.RawQuery == "code=errpass" {
		return nil, fmt.Errorf("invalid user/pass")
	}

	// Set session ID
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "PHPSESSID" {
			c.phpsessidCookie = cookie
		}
	}

	return c, nil
}
