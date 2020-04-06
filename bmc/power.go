package bmc

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// PowerState Different power states for the BMC
type PowerState string

// Power States
const (
	PowerStart PowerState = "poweron"
	PowerStop             = "poweroff"
	PowerReset            = "powerreboot"
)

// Start the server
func (c *Client) Start(ctx context.Context) error {
	return c.SetPowerState(ctx, PowerStart)
}

// Stop the server
func (c *Client) Stop(ctx context.Context) error {
	return c.SetPowerState(ctx, PowerStop)
}

// Reset the server
func (c *Client) Reset(ctx context.Context) error {
	return c.SetPowerState(ctx, PowerReset)
}

// SetPowerState Sets the server's power state
func (c *Client) SetPowerState(ctx context.Context, powerState PowerState) error {
	form := url.Values{}
	form.Add("power_option", string(powerState))
	req, err := c.buildRequest(ctx, "POST", fmt.Sprintf("https://%s/cgi_bin/ipmi_set_powercontrol.cgi", c.ip), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	return nil
}
