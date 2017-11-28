package main

import (
	"fmt"
	"sync"
	"time"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/coding-yogi/cfmonitoring/cf"
	"github.com/coding-yogi/cfmonitoring/influxdb"
	"github.com/coding-yogi/cfmonitoring/log"
)

func main() {
	c := cf.Client{}
	err := c.GetClient()
	failOnError(err, "Unable to connect to CF")

	apps := c.GetApps()

	ic := influxdb.Client{}
	err = ic.ConnectToInfluxDB()
	failOnError(err, "Unable to connect to influx DB")

	err = ic.CreateDatabase()
	failOnError(err, "Unable to create influxDB database")

	batchSize := 15
	var wg sync.WaitGroup

	for _, app := range apps {
		wg.Add(1)
		go func(app cfclient.App) {
			counter := 1

			bp, err := ic.CreateBatchPoints()
			failOnError(err, "Unable to create batchpoints")

			for {
				appStats, err := c.GetAppStats(app)
				failOnError(err, "Unable to get appstats")

				for k, stat := range appStats {
					u := stat.Stats.Usage
					log.Debug(fmt.Sprintf("Instance: %s CPU: %f RAM: %d Disk: %d", k, u.CPU, u.Mem, u.Disk))

					tags := map[string]string{"instance": app.Name + "_" + k}
					fields := map[string]interface{}{
						"cpu":    u.CPU,
						"memory": u.Mem,
						"disk":   u.Disk,
					}

					pt, err := ic.CreateNewPoint(app.Name, tags, fields)
					failOnError(err, "Unable to create new point")

					bp.AddPoint(pt)
					failOnError(err, "Unable to Setup Channel")
				}

				counter = counter + 1
				if counter%batchSize == 0 {
					log.Info("Writing to DB")
					ic.Write(bp)
				}

				//sleep
				time.Sleep(500 * time.Millisecond)
			}
		}(app)
	}

	wg.Wait()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
