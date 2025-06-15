package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"rpi-data-logger/database"
	"rpi-data-logger/model"
)

// Callback for processing incoming messages
func NewMessageHandler(db *database.DB) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		slog.Debug("Received MQTT message", "topic", topic)

		// Extract sensor ID from topic: "env/sensor/{sensor_id}/data"
		parts := strings.Split(topic, "/")
		if len(parts) != 4 {
			slog.Warn("Ignoring message on invalid topic", "topic", topic)
			return
		}
		sensorID := parts[2]

		// Unmarshal JSON payload
		var data model.SensorData
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			slog.Error("Failed to unmarshal JSON payload", "error", err, "sensor_id", sensorID)
			return
		}

		// Attempt to insert the data into the database
		err := db.InsertReading(sensorID, data)
		if err != nil {
			if errors.Is(err, database.ErrUnknownSensor) {
				slog.Warn(
					"Rejected data from unknown sensor",
					"sensor_id",
					sensorID,
				)
			} else {
				slog.Error("Failed to process sensor reading", "error", err, "sensor_id", sensorID)
			}
			return
		}

		slog.Info("Successfully logged data", "sensor_id", sensorID)
	}
}
