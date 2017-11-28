# cfmonitoring
Golang based tool for monitoring CPU and RAM of CF instances using Graphana and InfluxDB

## Usage
* Clone repo
* Make sure Graphana and InfluxDB is setup
* Update "cf" section in "config.json" as per your requirement
* Update "influxdb" section if influx DB is not running on localhost or at defined port
* Execute "go run main.go -cfusername=your_cf_username -cfpassword=your_cf_password"

Tool will connect to specified CF instances and push the CPU + RAM usage to influxDB

## Troubleshooting
contact aniket.gadre@sap.com
