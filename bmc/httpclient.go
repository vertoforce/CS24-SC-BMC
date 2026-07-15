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
			CipherSuites: []uint16{tls.TLS_RSA_WITH_RC4_128_SHA},
			// The BMC only speaks old TLS; Go 1.18+ defaults the client minimum
			// version to TLS 1.2, and RC4 cipher suites only exist in TLS 1.2
			// and below, so allow down to TLS 1.0
			MinVersion:         tls.VersionTLS10,
			InsecureSkipVerify: true,
		},
	}
	client = &http.Client{Transport: transport}

	// Get remote certificate and information
	dialer := &tls.Dialer{Config: transport.TLSClientConfig}
	netConn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, nil, 0, err
	}
	defer netConn.Close()
	conn := netConn.(*tls.Conn)
	certificates = conn.ConnectionState().PeerCertificates
	cipherSuite = conn.ConnectionState().CipherSuite

	return client, certificates, cipherSuite, nil
}

// buildRequest Builds the request with the appropriate headers and auth
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
