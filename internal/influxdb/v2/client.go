package v2

import (
	"context"

	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/unpoller/influxunifi/internal/influxdb/model"
)

const (
	defaultInfluxBucket = "unifi"
	defaultInfluxOrg    = "unifi"
	defaultInfluxURL    = "http://127.0.0.1:8086"
)

// Options defines the required options for the V2 client
type Options struct {
	URL    string
	Token  string
	Org    string
	Bucket string
}

// Client defines a V2 InfluxDB client
type Client struct {
	influx influx.Client
	bucket string
	org    string
}

func withDefaultOptions(options Options) Options {
	if options.URL == "" {
		options.URL = defaultInfluxURL
	}

	if options.Bucket == "" {
		options.Bucket = defaultInfluxBucket
	}

	if options.Org == "" {
		options.Org = defaultInfluxOrg
	}

	return options
}

// New creates a new client with the required options
func New(options Options) *Client {
	options = withDefaultOptions(options)
	return &Client{
		influx: influx.NewClient(options.URL, options.Token),
		bucket: options.Bucket,
		org:    options.Org,
	}
}

// Write writes metrics to InfluxDB
func (c *Client) Write(points ...*model.Point) error {
	writeAPI := c.influx.WriteAPIBlocking(c.org, c.bucket)
	influxPoints := make([]*write.Point, len(points))
	for i, point := range points {
		influxPoints[i] = influx.NewPoint(point.Measurement, point.Tags, point.Fields, point.Timestamp)
	}
	return writeAPI.WritePoint(context.TODO(), influxPoints...)
}
