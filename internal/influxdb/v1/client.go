package v1

import (
	"crypto/tls"
	"io/ioutil"
	"strings"

	influx "github.com/influxdata/influxdb1-client/v2"

	"github.com/wamphlett/influxunifi/internal/influxdb/model"
)

const (
	defaultInfluxDB       = "unifi"
	defaultInfluxUsername = "unifipoller"
	defaultInfluxURL      = "http://127.0.0.1:8086"
)

// Options defines the required options for the V1 client
type Options struct {
	URL       string
	Username  string
	Password  string
	DB        string
	VerifySSL bool
}

// Client defines a V1 InfluxDB client
type Client struct {
	influx influx.Client
	DB     string
}

type errorLogger func(msg string, v ...interface{})

func withDefaultOptions(options Options, logger errorLogger) Options {
	if options.URL == "" {
		options.URL = defaultInfluxURL
	}

	if options.Username == "" {
		options.Username = defaultInfluxUsername
	}

	if strings.HasPrefix(options.Password, "file://") {
		options.Password = getPassFromFile(strings.TrimPrefix(options.Password, "file://"), logger)
	}

	if options.Password == "" {
		options.Password = defaultInfluxUsername
	}

	if options.DB == "" {
		options.DB = defaultInfluxDB
	}

	return options
}

// New creates a new Client with the required options
func New(options Options, logger errorLogger) (*Client, error) {
	options = withDefaultOptions(options, logger)
	influx, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:      options.URL,
		Username:  options.Username,
		Password:  options.Password,
		TLSConfig: &tls.Config{InsecureSkipVerify: !options.VerifySSL}, // nolint: gosec
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		influx: influx,
		DB:     options.DB,
	}, nil
}

// Write writes metrics to InfluxDB
func (c *Client) Write(points ...*model.Point) error {
	batchPoints, err := influx.NewBatchPoints(influx.BatchPointsConfig{Database: c.DB})
	if err != nil {
		return err
	}
	for _, point := range points {
		point, err := influx.NewPoint(point.Measurement, point.Tags, point.Fields, point.Timestamp)
		if err != nil {
			continue
		}
		batchPoints.AddPoint(point)
	}
	return c.influx.Write(batchPoints)
}

func getPassFromFile(filename string, logger errorLogger) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		logger("Reading InfluxDB Password File: %v", err)
	}

	return strings.TrimSpace(string(b))
}
