package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerUseHTTPS bool
	ServerURL      string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewConfig() (*Config, error) {
	var serverURL string
	var err error
	serverUseHTTPS := flag.String("s", "false", "Use HTTPS connection to the server")
	serverAddr := flag.String("a", "localhost:8080", "Metrics server address (host:port)")
	pollInterval := flag.String("p", "2", "Interval for scraping metrics (in seconds)")
	reportInterval := flag.String("r", "10", "Interval for sending metrics (in seconds)")

	flag.Parse()
	if e := os.Getenv("USE_HTTPS"); e != "" {
		serverUseHTTPS = &e
	}
	if e := os.Getenv("ADDRESS"); e != "" {
		serverAddr = &e
	}
	if e := os.Getenv("POLL_INTERVAL"); e != "" {
		pollInterval = &e
	}
	if e := os.Getenv("REPORT_INTERVAL"); e != "" {
		reportInterval = &e
	}
	reportIntervalInt, err := strconv.Atoi(*reportInterval)
	if err != nil {
		return nil, err
	}
	pollIntervalInt, err := strconv.Atoi(*pollInterval)
	if err != nil {
		return nil, err
	}
	if pollIntervalInt > reportIntervalInt {
		return nil, fmt.Errorf("poll interval can not be more that report interval")
	}
	serverUseHTTPSBool, err := strconv.ParseBool(*serverUseHTTPS)
	if err != nil {
		return nil, err
	}
	if serverUseHTTPSBool {
		serverURL = fmt.Sprintf("https://%v/update", *serverAddr)
	} else {
		serverURL = fmt.Sprintf("http://%v/update", *serverAddr)
	}

	return &Config{
		ServerUseHTTPS: serverUseHTTPSBool,
		ServerURL:      serverURL,
		PollInterval:   time.Second * time.Duration(pollIntervalInt),
		ReportInterval: time.Second * time.Duration(reportIntervalInt),
	}, nil
}
