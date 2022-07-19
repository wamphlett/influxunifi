# influxunifi

## UnPoller Input Plugin

Collects UniFi data from a UniFi controller using the API.

# FORKED
Forked this in order to add support for the InfluxDB V2 client with Token authentication. V1 support remains the default and unchanged.

New Env vars
```env
UP_INFLUXDB_USEV2=bool # enables the V2 client, false by default
UP_INFLUXDB_TOKEN=string # authentication token
UP_INFLUXDB_ORG=string # Influx organisation
```
