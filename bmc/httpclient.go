package bmc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
)

// createHTTPClient Creates an HTTP client specifically made to communicate with the BMC
// The BMC uses a self-signed cert, and a very old cipher suite.  This client makes it so we can still communicate with it
// This function connects and grabs the server certificate and cipher suite, and builds a client set to communicate with that server
func createHTTPClient(ctx context.Context, ip string, port uint16) (client *http.Client, certificates []*x509.Certificate, cipherSuite uint16, err error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			// BMC uses cipher TLS_RSA_WITH_RC4_128_SHA
			// See this list: https://www.iana.org/assignments/tls-parameters/tls-parameters.xml
			CipherSuites:       []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
			InsecureSkipVerify: true,
		},
	}
	client = &http.Client{Transport: transport}

	// Get remote certificate and information
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), transport.TLSClientConfig)
	if err != nil {
		return nil, nil, 0, err
	}
	certificates = conn.ConnectionState().PeerCertificates
	cipherSuite = conn.ConnectionState().CipherSuite

	return client, certificates, cipherSuite, nil
}

// buildRequeste Builds the request with the appropriate headers and auth
func (c *Client) buildRequest(ctx context.Context, method string, URL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, URL, body)
	if err != nil {
		return nil, err
	}

	// Add content type if we have a body
	if body != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	// Add auth header if we have one
	if c.phpsessidCookie != nil {
		req.AddCookie(c.phpsessidCookie)
	}

	return req, nil
}
