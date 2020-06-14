package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/namsral/flag"
	"github.com/vertoforce/CS24-SC-BMC/bmc"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Flags
var (
	IP        string
	Username  string
	Password  string
	Port      uint
	Action    string
	LogFormat string
)

func main() {
	log := logrus.New()

	// Parse flags
	flag.StringVar(&IP, "IP", "", "IP of server to connect to")
	flag.UintVar(&Port, "Port", 443, "Port of server to connect to")
	flag.StringVar(&Action, "Action", "info", "Action to perform on server. Options are: info, start, stop, reset, monitor")
	flag.StringVar(&Username, "Username", "", "Username for BMC")
	flag.StringVar(&Password, "Password", "", "Password for BMC")
	flag.StringVar(&LogFormat, "LogFormat", "text", "The formatting of the logs, can be text or json.")
	flag.Parse()

	switch LogFormat {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	// Watch for program cancelation
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer func() {
		signal.Stop(signals)
		cancel()
	}()
	go func() {
		select {
		case <-signals:
			log.Error("Canceling")
			cancel()
		case <-ctx.Done():
		}
	}()

	// Connect to server
	log.WithFields(logrus.Fields{"IP": IP, "Port": Port}).Info("Connecting to server")
	c, err := bmc.New(ctx, IP, uint16(Port), Username, Password)
	if err != nil {
		log.WithError(err).Error("Failed to connect to BMC server")
		return
	}

	switch Action {
	case "start":
		log.Info("Starting server")
		err = c.Start(ctx)
		if err != nil {
			log.WithError(err).Error("Error")
		}
	case "stop":
		log.Info("Stopping server")
		err = c.Stop(ctx)
		if err != nil {
			log.WithError(err).Error("Error")
		}
	case "reset":
		log.Info("Resetting server")
		err = c.Reset(ctx)
		if err != nil {
			log.WithError(err).Error("Error")
		}
	case "info":
		certificateSubjects := []string{}
		for _, certificate := range c.Certificates {
			certificateSubjects = append(certificateSubjects, certificate.Subject.String())
		}
		temperatures, _ := c.GetTemperature(ctx)
		log.WithFields(logrus.Fields{
			"CertificateSubjects": certificateSubjects,
			"CipherSuiteCode":     c.CipherSuite,
			"Temperatures":        temperatures,
		}).Info("Server Info")
	case "monitor":
		for {
			monitor(ctx, log)
		}
	}

	log.Info("Done")
}

// A single monitor run
func monitor(ctx context.Context, log *logrus.Logger) {
	defer time.Sleep(time.Second * 30)
	// Connect
	c, err := bmc.New(ctx, IP, uint16(Port), Username, Password)
	if err != nil {
		log.WithError(err).Error("Could not connect to server")
		return
	}

	temps, err := c.GetTemperature(ctx)
	log.WithFields(logrus.Fields{
		"Temperature": temps,
	}).Info("Server Info")
}
