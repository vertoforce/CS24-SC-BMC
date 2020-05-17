package bmc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	temperatureEndpoint = "https://%s:%d/cgi_bin/ipmi_get_info.cgi?operation=temperature"
)

// TempReading is an individual reading from a sensor
type TempReading struct {
	Description   string
	SensorReading float64
}

// TempReadings contains physical sensor temperature readings for the server CPUs
type TempReadings []TempReading

// GetTemperature Gets the temperature readings of the server
func (c *Client) GetTemperature(ctx context.Context) (TempReadings, error) {
	req, err := c.buildRequest(ctx, "GET", fmt.Sprintf(temperatureEndpoint, c.ip, c.port), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := TempReadings{}

	// Go through every table
	doc.Find("table:contains(Description)[border='1']").Each(func(i int, s *goquery.Selection) {
		var description string
		var sensorReading float64

		// For each row, find the values we care about
		s.Find("tr").Each(func(t int, tr *goquery.Selection) {
			tds := tr.Find("td")
			if tds.Length() < 2 {
				return
			}
			// Get label and value
			label := strings.ToLower(strings.Trim(tds.First().Text(), " \t\n"))
			value := strings.Trim(tds.First().Next().Text(), " \t\n")
			switch label {
			case "description:":
				description = value
			case "sensorreading:":
				sensorReading, err = strconv.ParseFloat(value, 10)
			}
		})

		ret = append(ret, TempReading{
			Description:   description,
			SensorReading: sensorReading,
		})
	})

	return ret, nil
}
