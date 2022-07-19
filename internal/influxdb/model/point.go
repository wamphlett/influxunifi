package model

import "time"

// Point defines the fields required in order to publish a metric to InfluxDB
type Point struct {
	Measurement string
	Fields      map[string]interface{}
	Tags        map[string]string
	Timestamp   time.Time
}
