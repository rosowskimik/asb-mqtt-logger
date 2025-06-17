# ASB MQTT Logger

## General Description

This application is a Go service designed to run on a Linux system.
It acts as the central data collector for a network of environment sensors.

The service performs the following functions:

- Connects to a local MQTT broker.
- Subscribes to the topic `env/sensor/+/data` to listen for messages from all registered sensors.
- Parses incoming JSON data containing temperature, pressure, and humidity statistics.
- Saves valid sensor readings to a local SQLite database.

The database is structured with one `sensor_readings` table for the time-series data.

## Build Instructions

A Go toolchain (version 1.24.x or newer) is required.

To build the application, navigate to the `abs-mqtt-logger` directory and run:

```bash
go build .
```

This will produce a single executable file named `rpi-data-logger`.

## Usage Instructions

1. **Run the Service:**
   Execute the binary from your terminal:

   ```bash
   ./rpi-data-logger
   ```

   For continuous operation, it's recommended to run this as a `systemd` service.

1. **Configuration:**
   The application is configured via environment variables:

   - `APP_DB_PATH`: The file path for the SQLite database. (Default: `environment_data.db`)
   - `APP_MQTT_BROKER`: The hostname or IP of the MQTT broker. (Default: `localhost`)
   - `APP_MQTT_PORT`: The port for the MQTT broker. (Default: `1883`)
   - `APP_LOG_LEVEL`: The logging level. (Options: `DEBUG`, `INFO`, `WARN`, `ERROR`. Default: `INFO`)

## License

This project is licensed under the Apache 2.0 License. See the `LICENSE` file in the root of the repository for more details.
