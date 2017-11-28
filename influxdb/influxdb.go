package influxdb

import (
	"fmt"
	"time"

	"github.com/coding-yogi/cfmonitoring/config"
	"github.com/coding-yogi/cfmonitoring/log"
	client "github.com/influxdata/influxdb/client/v2"
)

type Client struct {
	client client.Client
	db     string
}

// ConnectToInfluxDB ...
func (c *Client) ConnectToInfluxDB() error {
	settings, err := config.GetConfig()
	if err != nil {
		return err
	}

	config := client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s:%d", settings.InfluxDB.Host, settings.InfluxDB.Port),
	}

	log.Info("Connecting to influxDB at " + config.Addr)
	c.client, err = client.NewHTTPClient(config)
	if err != nil {
		return err
	}

	c.db = settings.InfluxDB.Database
	return nil
}

// Create DB
func (c *Client) CreateDatabase() error {
	q := client.NewQuery("CREATE DATABASE "+c.db, "", "")
	if _, err := c.client.Query(q); err != nil {
		return err
	}

	return nil
}

// CreateBatchPoints ...
func (c *Client) CreateBatchPoints() (client.BatchPoints, error) {
	return client.NewBatchPoints(client.BatchPointsConfig{
		Database:  c.db,
		Precision: "s",
	})
}

// CreateNewPoint ..
func (c *Client) CreateNewPoint(name string, tags map[string]string, fields map[string]interface{}) (*client.Point, error) {
	return client.NewPoint(name, tags, fields, time.Now())
}

// Write
func (c *Client) Write(bp client.BatchPoints) error {
	return c.client.Write(bp)
}
